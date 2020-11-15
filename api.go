package main

import (
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func getTodo(c *gin.Context) {

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res := db.Client.Keys(db.Ctx, "*")

	keys, err := res.Result()

	if err != nil {
		c.String(500, "can't get keys")
	} else {
		r := strings.Split("", "")
		for _, v := range keys {
			r = append(r, v, "\n")
			value, _ := db.Client.Get(db.Ctx, v).Result()
			r = append(r, value, "\n")

		}
		c.String(200, strings.Join(r, ""))
	}

	defer db.dbClose()
}

func postTodo(c *gin.Context) {

	key := c.PostForm("key")

	value := c.PostForm("value")

	if key == "" && value == "" {
		c.String(400, "wrong request")
		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res := db.Client.Set(db.Ctx, key, value, 0)

	c.String(200, res.String())

	defer db.dbClose()

}

func delTodo(c *gin.Context) {

	name := c.Param("name")

	if name == "" {
		c.String(400, "wrong param")
		return
	}

	db := db{
		Ctx:    nil,
		Client: nil,
	}

	db.dbInit()

	res := db.Client.Del(db.Ctx, name)

	_, err := res.Result()
	if err != nil {
		c.String(500, "del error")
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
}

func postDownload(c *gin.Context) {
}

func delDownload(c *gin.Context) {
}

func getShortUrl(c *gin.Context) {

}

func postShortUrl(c *gin.Context) {

}

func delShortUrl(c *gin.Context) {

}
