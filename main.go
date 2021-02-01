package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
)

var debug = os.Getenv("cbox_debug")

func main() {

	var cerPath string
	var keyPath string
	var addr string
	var logPath string

	flag.StringVar(
		&keyPath, "k", "/etc/nginx/notok.cf.key", "tls key path")

	flag.StringVar(
		&cerPath, "c", "/etc/nginx/notok.cf.cer", "tls cer path")

	flag.StringVar(
		&addr, "a", ":55557", "port")

	flag.StringVar(
		&logPath, "d", "./", "log path")

	f, _ := os.Create(logPath + "/cbox.log")
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
		_ = r.RunTLS(addr, cerPath, keyPath)

	}

}
