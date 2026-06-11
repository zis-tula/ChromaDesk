<template>
  <div>
    <h2 class="page-title">{{ $t('about.title') }}</h2>

    <div class="card">
      <div class="about-logo">
        <div class="about-logo-icon"><img src="/images/logo.png" alt="ChromaDesk" /></div>
        <div>
          <div class="about-app-name">ChromaDesk</div>
          <div class="about-version">{{ $t('version.label') }}</div>
        </div>
      </div>
      <div class="separator"></div>
      <p class="about-desc">{{ $t('about.desc') }}</p>
    </div>

    <div class="card">
      <div class="card-title">{{ $t('about.links') }}</div>
      <a class="link-item" href="https://github.com/barry-ran/ChromaDesk" target="_blank" rel="noopener">
        <div class="link-icon">🌐</div>
        <div class="link-info">
          <div class="link-title">{{ $t('about.sourceCode') }}</div>
          <div class="link-url">github.com/barry-ran/ChromaDesk</div>
        </div>
        <span class="link-arrow">→</span>
      </a>
      <a class="link-item" href="https://github.com/barry-ran/chromadesk-remoting" target="_blank" rel="noopener">
        <div class="link-icon">🔗</div>
        <div class="link-info">
          <div class="link-title">{{ $t('about.basedOn') }}</div>
          <div class="link-url">github.com/barry-ran/chromadesk-remoting</div>
        </div>
        <span class="link-arrow">→</span>
      </a>
    </div>

    <div class="card">
      <div class="card-title">{{ $t('about.license') }}</div>
      <p class="about-desc" style="margin-bottom:0">{{ $t('about.licenseDesc') }}</p>
    </div>
  </div>
</template>

<style scoped>
.about-logo { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; }
.about-logo-icon { width: 56px; height: 56px; border-radius: 12px; overflow: hidden; }
.about-logo-icon img { width: 100%; height: 100%; object-fit: cover; }
.about-app-name { font-size: 22px; font-weight: 700; }
.about-version { font-size: 14px; color: var(--text-secondary); margin-top: 2px; }
.about-desc { font-size: 14px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 16px; }
.separator { height: 1px; background: var(--border); margin: 16px 0; }

.link-item { display: flex; align-items: center; gap: 12px; padding: 12px; border-radius: var(--radius); cursor: pointer; transition: background 0.15s; text-decoration: none; color: var(--text-primary); }
.link-item:hover { background: rgba(255,255,255,0.05); }
.link-icon { width: 32px; height: 32px; border-radius: 8px; background: rgba(233,69,96,0.15); display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; }
.link-info { flex: 1; }
.link-title { font-size: 14px; font-weight: 500; }
.link-url { font-size: 12px; color: var(--text-secondary); margin-top: 2px; }
.link-arrow { color: var(--text-secondary); font-size: 14px; }
</style>
