package swagger

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openAPIFS embed.FS

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/swagger", h.UI)
	r.GET("/swagger/openapi.yaml", h.Spec)
}

func (h *Handler) Spec(c *gin.Context) {
	bytes, err := openAPIFS.ReadFile("openapi.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load openapi spec"})
		return
	}
	c.Data(http.StatusOK, "application/yaml", bytes)
}

func (h *Handler) UI(c *gin.Context) {
	html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <title>Student & T API Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/swagger/openapi.yaml',
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis],
    });
  </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
