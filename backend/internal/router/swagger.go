package router

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed swagger.json
var swaggerJSON []byte

// swaggerHTML 基于 Swagger UI CDN 的文档页（离线时回退显示原始 JSON）。
const swaggerHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8"/>
<title>处方智能审核引擎 API 文档</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = function(){
  try {
    SwaggerUIBundle({
      url: "/swagger/doc.json",
      dom_id: "#swagger-ui",
      deepLinking: true
    });
  } catch(e) {
    document.getElementById("swagger-ui").innerHTML = "<pre>" + "Swagger UI 加载失败，请访问 /swagger/doc.json 查看 OpenAPI 定义。" + "</pre>";
  }
};
</script>
</body>
</html>`

// registerSwaggerRoutes 注册 Swagger 文档路由。
func registerSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger/doc.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", swaggerJSON)
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
