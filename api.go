package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/cavaliercoder/grab"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getServerIp(c *gin.Context) {

	t := "https://httpbin.org/ip"

	response, err := http.Get(t)

	if err != nil {
		c.String(500, "get server ip error")
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.String(500, "get server ip then read data error")
		return
	}

	val := string(data)

	res := gjson.Get(val, "origin").String()

	c.String(200, res)
}

func getClientIp(c *gin.Context) {
	clientIp := c.ClientIP()

	c.String(200, clientIp)
}

func getWeather(c *gin.Context) {

	//city := c.Param("city")

	u, err := url.Parse("https://weather.com/zh-CN/weather/today/l/76d8b2d61e3b49327ef17615282eb460cfaf4b511a2ff907b5caae4e66badb0c")

	if err != nil {
		c.String(500, "parse url error")
		return
	}

	t := u.String()

	doc, err := htmlquery.LoadURL(t)
	if err != nil {
		c.String(500, "htmlquery error")
		return
	}

	xpath := "/html/body/div[1]/main/div[2]/div[2]/div[1]/div/section/div/div[2]/div[1]/span"

	res := htmlquery.FindOne(doc, xpath)

	weather := htmlquery.InnerText(res)

	c.String(200, weather)

}

func getToDoApiAddr(c *gin.Context) {

	httpUrl := "http://notok.cf:55557/api/todo"
	httpsUrl := "https://notok.cf:55557/api/todo"

	if debug == "1" {
		c.String(200, httpUrl)
	} else {
		c.String(200, httpsUrl)
	}

}

func getTodo(c *gin.Context) {

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res := db.Client.LRange(db.Ctx, "todo", 0, -1)

	keys, err := res.Result()

	if err != nil {
		c.String(500, "getTodo get range error ")
		return
	} else {
		r := ""
		for _, v := range keys {

			r, _ = sjson.Set(r, "title.-1", v)
			value, _ := db.Client.Get(db.Ctx, v).Result()
			r, _ = sjson.Set(r, "data.-1", value)

		}

		c.String(200, r)
	}

	defer db.dbClose()
}

func getOneTodo(c *gin.Context) {

	name := c.Param("name")

	fmt.Println(name)

	if name == "" {
		c.String(400, "getOneTodo get param error")
		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res, err := db.Client.Get(db.Ctx, name).Result()

	if err != redis.Nil {
		c.Redirect(http.StatusMovedPermanently, res)
	} else {
		c.String(500, "getOneTodo get value error ")
		return
	}

	defer db.dbClose()

}

func postTodo(c *gin.Context) {

	key := c.PostForm("key")

	value := c.PostForm("value")

	if key == "" && value == "" {
		c.String(400, "postTodo get param error")
		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	request := db.Client.Get(db.Ctx, key)

	res, err := request.Result()

	if err != nil && err.Error() != "redis: nil" {
		c.String(500, "posTodo set string error:%s", err.Error())
		return
	}

	if res == "" {

		pushRes := db.Client.RPush(db.Ctx, "todo", key)

		if _, err := pushRes.Result(); err != nil {

			_ = db.Client.Del(db.Ctx, key) // error rollback

			c.String(500, "posTodo push string error")
			return
		}
	}

	setRes := db.Client.Set(db.Ctx, key, value, 0)

	if _, err := setRes.Result(); err != nil {

		c.String(500, "posTodo set string error")
		return
	}

	defer db.dbClose()

	c.String(200, "posTodo ok")
}

func delTodo(c *gin.Context) {

	name := c.Param("name")

	if name == "" {
		c.String(400, "delTodo get param error")
		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	delRes := db.Client.Del(db.Ctx, name)

	remRes := db.Client.LRem(db.Ctx, "todo", 0, name)

	_, err := delRes.Result()
	_, err2 := remRes.Result()

	if err != nil || err2 != nil {
		c.String(500, "del error")
		return
	} else {
		c.String(200, "del success")
	}

	defer db.dbClose()

}

func getIthomeNew(c *gin.Context) {

	jsonReturn := ""
	ua := "Mozilla/5.0 (Linux; Android 7.1.1; OPPO R9sk) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/76.0.3809.111 Mobile Safari/537.36\""
	var coll = colly.NewCollector(
		colly.UserAgent(ua),
	)

	coll.OnRequest(func(r *colly.Request) {

		//fmt.Println(r.URL)
	})

	coll.OnError(
		func(response *colly.Response, err error) {
			log.Fatal(err)
		})

	coll.OnHTML("body > div.index-box > div.content",
		func(c *colly.HTMLElement) {

			c.ForEach(
				"div > a > div.plc-image > img ",
				func(i int, element *colly.HTMLElement) {
					v := element.Attr("data-original")
					jsonReturn, _ = sjson.Set(jsonReturn, "imagesUrl.-1", v)

					//fmt.Println()
				})
			c.ForEach("div > a",
				func(i int, element *colly.HTMLElement) {
					v := element.Attr("href")
					jsonReturn, _ = sjson.Set(jsonReturn, "newsUrl.-1", v)

					contextColl := colly.NewCollector(
						colly.UserAgent(ua),
					)

					contextNews := ""

					contextColl.OnHTML("body > div.page-box > div.con-box > div.news > main > div.news-content",
						func(e *colly.HTMLElement) {
							e.ForEach("p",
								func(i int, element *colly.HTMLElement) {

									if element.Text == "" {
										contextNews += element.ChildAttr("img", "data-original")

									}
									contextNews += element.Text + "\n"
								})
						},
					)
					_ = contextColl.Visit(v)

					jsonReturn, _ = sjson.Set(jsonReturn, "context.-1", contextNews)
				})

			c.ForEach(
				"div > a > div.plc-con > p ",
				func(i int, element *colly.HTMLElement) {
					v := element.Text
					jsonReturn, _ = sjson.Set(jsonReturn, "title.-1", v)

				})

		},
	)

	_ = coll.Visit("https://m.ithome.com")

	c.String(200, jsonReturn)

}

func postIthomeNew(c *gin.Context) {

}

func getDownload(c *gin.Context) {
	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res, err := db.Client.LRange(db.Ctx, "download", 0, -1).Result()

	if err != nil {
		c.String(500, "getDownload get data error : %s", err.Error())
	}

	jsonReturn := ""

	for _, v := range res {
		jsonReturn, _ = sjson.Set(jsonReturn, "download.-1", v)

		res := db.Client.HMGet(db.Ctx, v, "url", "progress", "done")

		for i, val := range res.Val() {
			switch i {
			case 0:
				jsonReturn, _ = sjson.Set(jsonReturn, "url.-1", val)
			case 1:
				jsonReturn, _ = sjson.Set(jsonReturn, "progress.-1", val)
			case 2:
				jsonReturn, _ = sjson.Set(jsonReturn, "done.-1", val)

			}
		}

	}

	c.String(200, jsonReturn)
}

func postDownload(c *gin.Context) {
	target := c.PostForm("target")

	if target == "" {
		c.String(400, "postDownload get param error")
		return
	}

	client := grab.NewClient()

	request, _ := grab.NewRequest(".", target)

	response := client.Do(request)

	if response.Err() != nil {
		c.String(500, "postDownload start download error:"+
			response.Err().Error())
		return
	}

	filename := response.Filename

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	db.Client.RPush(db.Ctx, "download", filename)

	db.Client.HMSet(db.Ctx, filename, "url", target, "progress", "0",
		"done", "false")

	defer db.dbClose()

	c.String(200, "start download , file name is %s", filename)
}

func delDownload(c *gin.Context) {

	filename := c.Param("filename")

	if filename == "" {
		c.String(400, "delDownload get param error")

		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	db.Client.LRem(db.Ctx, "download", 0, filename)
	db.Client.HDel(db.Ctx, filename, "url", "progress", "done")

	err := os.Remove(filename)

	if err != nil {
		c.String(500, "delDownload del file error:%s", err.Error())
		return
	}

	c.String(200, "delDownload success")

}

func getDu(c *gin.Context) {
	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	defer db.dbClose()

	num := rand.Uint64() % 1219

	res := db.Client.LRange(db.Ctx, "du", int64(num), int64(num))

	val, err := res.Result()
	if err != nil {
		c.String(500, "getDu error:%s", err.Error())
		return
	}

	c.String(200, strings.Join(val, ""))
}
