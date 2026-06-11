package handler

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterAdminUI serves the embedded Vue admin frontend at /admin/.
func RegisterAdminUI(router *gin.Engine, webFS embed.FS) {
	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		return
	}

	// Check if dist has real content (not just placeholder)
	entries, _ := fs.ReadDir(distFS, ".")
	hasIndex := false
	for _, e := range entries {
		if e.Name() == "index.html" {
			hasIndex = true
			break
		}
	}

	if !hasIndex {
		// Vue not built yet — serve a helpful message
		router.GET("/admin", func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8",
				[]byte("<h2>Admin UI not built</h2><p>Run <code>cd web && npm install && npm run build</code> then restart the server.</p>"))
		})
		return
	}

	fileServer := http.StripPrefix("/admin/", http.FileServer(http.FS(distFS)))

	router.GET("/admin/*filepath", func(c *gin.Context) {
		path := strings.TrimPrefix(c.Param("filepath"), "/")

		// Serve exact file if it exists
		if path != "" {
			if _, err := fs.Stat(distFS, path); err == nil {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// SPA fallback: serve index.html
		indexData, _ := fs.ReadFile(distFS, "index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	})

	router.GET("/admin", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/admin/")
	})
}
