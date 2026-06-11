// Copyright 2026 ChromaDesk Authors
// SkillRegistry — scans the skills directory, parses SKILL.md frontmatter,
// and starts MCP server subprocesses for each applicable skill.

use serde::Deserialize;
use serde_json::Value;
use std::collections::HashMap;
use std::path::{Path, PathBuf};
use tracing::{info, warn};

use crate::dep_checker;
use crate::mcp_client::McpClient;

/// Parsed frontmatter from a SKILL.md file.
#[derive(Debug, Deserialize, Default)]
#[allow(dead_code)]
pub struct SkillMeta {
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub metadata: SkillMetadata,
}

#[derive(Debug, Deserialize, Default)]
pub struct SkillMetadata {
    #[serde(default)]
    pub openclaw: Option<OpenClawMeta>,
}

#[derive(Debug, Deserialize, Default, Clone)]
pub struct OpenClawMeta {
    #[serde(default)]
    pub os: Vec<String>,
    #[serde(default)]
    pub requires: OpenClawRequires,
    #[serde(default)]
    pub install: Vec<InstallStep>,
}

#[derive(Debug, Deserialize, Default, Clone)]
pub struct OpenClawRequires {
    #[serde(default)]
    pub bins: Vec<String>,
}

#[derive(Debug, Deserialize, Default, Clone)]
#[allow(dead_code)]
pub struct InstallStep {
    pub id: String,
    pub kind: String,
    #[serde(default)]
    pub package: String,
}

/// A loaded skill, ready to serve tool calls.
#[allow(dead_code)]
pub struct LoadedSkill {
    pub meta: SkillMeta,
    pub client: McpClient,
}

/// Describes why a skill failed to load.
#[derive(Debug)]
pub struct SkillLoadError {
    pub skill_name: String,
    pub reason: String,
    pub missing_bins: Vec<String>,
}

pub struct SkillRegistry {
    skills_dirs: Vec<PathBuf>,
    skills: HashMap<String, LoadedSkill>,
    load_errors: Vec<SkillLoadError>,
}

impl SkillRegistry {
    pub fn new(skills_dirs: &[String]) -> Self {
        Self {
            skills_dirs: skills_dirs.iter().map(PathBuf::from).collect(),
            skills: HashMap::new(),
            load_errors: Vec::new(),
        }
    }

    /// Scan all skills directories and start each applicable skill's MCP server.
    pub async fn load(&mut self) {
        for dir in self.skills_dirs.clone() {
            if !dir.exists() {
                warn!("skills directory not found: {}", dir.display());
                continue;
            }

            let entries = match std::fs::read_dir(&dir) {
                Ok(e) => e,
                Err(e) => {
                    warn!("cannot read skills directory {}: {}", dir.display(), e);
                    continue;
                }
            };

            for entry in entries.flatten() {
                let skill_dir = entry.path();
                if !skill_dir.is_dir() {
                    continue;
                }
                let skill_md = skill_dir.join("SKILL.md");
                if !skill_md.exists() {
                    continue;
                }

                if self.skills.contains_key(&skill_dir.file_name()
                    .unwrap_or_default().to_string_lossy().to_string())
                {
                    continue;
                }

                match self.load_skill(&skill_dir, &skill_md).await {
                    Ok(name) => info!("Loaded skill: {}", name),
                    Err(e) => warn!(
                        "Failed to load skill at {}: {}",
                        skill_dir.display(),
                        e
                    ),
                }
            }
        }
    }

    /// Reload a single skill by name (used after user installs dependencies).
    pub async fn reload_skill(&mut self, skill_name: &str) -> Result<(), String> {
        self.skills.remove(skill_name);
        self.load_errors.retain(|e| e.skill_name != skill_name);

        // Search all configured directories
        for dir in &self.skills_dirs {
            let skill_dir = dir.join(skill_name);
            let skill_md = skill_dir.join("SKILL.md");
            if skill_md.exists() {
                return self.load_skill(&skill_dir, &skill_md).await
                    .map(|_| ())
                    .map_err(|e| e.to_string());
            }
        }

        // Fallback: user skills directory
        let user_dir = get_user_skills_dir();
        let user_skill_dir = user_dir.join(skill_name);
        let user_skill_md = user_skill_dir.join("SKILL.md");
        if user_skill_md.exists() {
            return self.load_skill(&user_skill_dir, &user_skill_md).await
                .map(|_| ())
                .map_err(|e| e.to_string());
        }

        Err(format!("skill '{}' not found", skill_name))
    }

    /// Return a flat list of all tools from all running skill servers.
    pub fn list_tools(&self) -> Vec<Value> {
        self.skills
            .values()
            .flat_map(|s| s.client.cached_tools())
            .collect()
    }

    /// Return all load errors from the last load.
    pub fn load_errors(&self) -> &[SkillLoadError] {
        &self.load_errors
    }

    /// Call a tool by name. Finds the owning skill and delegates.
    pub async fn call_tool(
        &self,
        tool_name: &str,
        args: Value,
    ) -> anyhow::Result<Value> {
        for skill in self.skills.values() {
            if skill.client.has_tool(tool_name) {
                return skill.client.call_tool(tool_name, args).await;
            }
        }
        anyhow::bail!("unknown tool: {}", tool_name)
    }

    // ---- private ----

    async fn load_skill(
        &mut self,
        skill_dir: &Path,
        skill_md: &Path,
    ) -> anyhow::Result<String> {
        let content = std::fs::read_to_string(skill_md)?;
        let meta = parse_frontmatter(&content)?;
        let skill_name = meta.name.clone();

        // Check OS compatibility
        if let Some(ref oc) = meta.metadata.openclaw {
            if !oc.os.is_empty() {
                let current_os = std::env::consts::OS;
                let oc_os = match current_os {
                    "windows" => "win32",
                    "macos" => "darwin",
                    other => other,
                };
                if !oc.os.iter().any(|o| o == oc_os) {
                    info!(
                        "Skill '{}' does not support OS '{}', skipping",
                        meta.name, oc_os
                    );
                    return Ok(meta.name);
                }
            }
        }

        // Determine install kind and handle accordingly
        let oc = meta.metadata.openclaw.as_ref()
            .ok_or_else(|| anyhow::anyhow!("no openclaw metadata"))?;

        let step = oc.install.first()
            .ok_or_else(|| anyhow::anyhow!("no install step"))?;

        let is_binary = step.kind == "binary";

        // For non-binary skills, check requires.bins
        if !is_binary && !oc.requires.bins.is_empty() {
            let missing = dep_checker::check_bins(&oc.requires.bins);
            if !missing.is_empty() {
                let err = SkillLoadError {
                    skill_name: skill_name.clone(),
                    reason: "missing_deps".to_string(),
                    missing_bins: missing.clone(),
                };
                warn!(
                    "Skill '{}' missing dependencies: {:?}",
                    skill_name, missing
                );
                self.load_errors.push(err);
                return Ok(skill_name);
            }
        }

        // Build the start command
        let cmd = build_start_command(step, skill_dir)?;

        let mut client = McpClient::new(&meta.name, cmd);
        match client.start().await {
            Ok(()) => {}
            Err(e) => {
                let err = SkillLoadError {
                    skill_name: skill_name.clone(),
                    reason: "start_failed".to_string(),
                    missing_bins: Vec::new(),
                };
                self.load_errors.push(err);
                return Err(e);
            }
        }
        client.fetch_tools().await?;

        self.skills.insert(skill_name.clone(), LoadedSkill { meta, client });
        Ok(skill_name)
    }
}

/// Parse YAML frontmatter from a SKILL.md file.
fn parse_frontmatter(content: &str) -> anyhow::Result<SkillMeta> {
    let mut lines = content.lines();
    if lines.next().map(str::trim) != Some("---") {
        anyhow::bail!("SKILL.md does not start with ---");
    }

    let mut yaml = String::new();
    for line in lines {
        if line.trim() == "---" {
            break;
        }
        yaml.push_str(line);
        yaml.push('\n');
    }

    let meta: SkillMeta = serde_yaml::from_str(&yaml)
        .map_err(|e| anyhow::anyhow!("YAML parse error: {}", e))?;
    Ok(meta)
}

/// Resolve the path to a binary skill executable.
/// Search order:
///   1. `<skill_dir>/<package>(.exe)` — the skill's own subdirectory
///   2. `<skill_host_exe_dir>/skills/<package>/<package>(.exe)` — installed layout
///   3. `<skill_host_exe_dir>/skills/<package>(.exe)` — legacy flat layout
///   4. `<skill_host_exe_dir>/<package>(.exe)` — fallback
fn resolve_binary_path(package: &str, skill_dir: &Path) -> anyhow::Result<PathBuf> {
    let bin_name = if cfg!(target_os = "windows") {
        format!("{}.exe", package)
    } else {
        package.to_string()
    };

    // 1. Look inside the skill directory (covers both built-in and user-added skills)
    let in_skill = skill_dir.join(&bin_name);
    if in_skill.exists() {
        return Ok(in_skill);
    }

    // 2-4. Fallback via skill-host executable directory
    if let Ok(exe) = std::env::current_exe() {
        if let Some(exe_dir) = exe.parent() {
            let sub = exe_dir.join("skills").join(package).join(&bin_name);
            if sub.exists() {
                return Ok(sub);
            }
            let flat = exe_dir.join("skills").join(&bin_name);
            if flat.exists() {
                return Ok(flat);
            }
            let beside = exe_dir.join(&bin_name);
            if beside.exists() {
                return Ok(beside);
            }
        }
    }

    anyhow::bail!(
        "binary '{}' not found in '{}' or skill-host directory",
        bin_name,
        skill_dir.display()
    )
}

/// Build the command to start the MCP server from the install step.
fn build_start_command(
    step: &InstallStep,
    skill_dir: &Path,
) -> anyhow::Result<Vec<String>> {
    match step.kind.as_str() {
        "binary" => {
            let bin_path = resolve_binary_path(&step.package, skill_dir)?;
            Ok(vec![bin_path.to_string_lossy().to_string()])
        }
        "npm" => {
            let pkg = &step.package;
            Ok(vec!["npx".to_string(), pkg.clone(), "--stdio".to_string()])
        }
        "uvx" => {
            let pkg = &step.package;
            Ok(vec!["uvx".to_string(), pkg.clone()])
        }
        other => {
            warn!("unknown install kind '{}', trying as command", other);
            Ok(vec![step.package.clone()])
        }
    }
}

fn get_user_skills_dir() -> PathBuf {
    if let Some(home) = dirs_home() {
        home.join(".chromadesk-skill-host").join("skills")
    } else {
        PathBuf::from(".chromadesk-skill-host/skills")
    }
}

fn dirs_home() -> Option<PathBuf> {
    #[cfg(target_os = "windows")]
    {
        std::env::var("USERPROFILE").ok().map(PathBuf::from)
    }
    #[cfg(not(target_os = "windows"))]
    {
        std::env::var("HOME").ok().map(PathBuf::from)
    }
}
