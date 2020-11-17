package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var debug = os.Getenv("cbox_debug")

func main() {

	f, _ := os.Create("cbox.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	// _cat api group
	_cat := r.Group("/_cat")
	{
		_cat.GET("/serverIp", getServerIp)

		_cat.GET("/clientIp", getClientIp)

		_cat.GET("/weather/:city", getWeather)

		_cat.GET("/toDoApiAddr", getToDoApiAddr)

		_cat.GET("/du", getDu)
	}

	// do api group
	do := r.Group("/do")
	{
		do.GET("/led", func(c *gin.Context) {

		})
	}

	// common api group
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {

			c.String(200, "pong")
		})

		api.GET("/ithomeNews", getIthomeNew)

		api.POST("/ithomeNews", postIthomeNew)

		api.GET("/toDo", getTodo)

		api.POST("/toDo", postTodo)

		api.DELETE("/toDo/:name", delTodo)

		api.GET("/shortUrl", getShortUrl)

		api.POST("/shortUrl", postShortUrl)

		api.DELETE("/shortUrl/:url", delShortUrl)

		api.GET("download", getDownload)

		api.POST("/download", postDownload)

		api.DELETE("download/:filename", delDownload)

	}

	if debug == "1" {
		_ = r.Run(":55557")
	} else {
		// running at tls
		_ = r.RunTLS(":55557", "/etc/nginx/notok.cf.cer",
			"/etc/nginx/notok.cf.key")

	}

}
