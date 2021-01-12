package main

import (
	"GoBlog/config"
	_ "GoBlog/config/nsq"
	"GoBlog/config/redis"
	"GoBlog/lib/mysql"
	"GoBlog/logic/controller"
	"GoBlog/service"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	mysql.ConnectMysql(config.MySQL, "default")
	defer mysql.DisconnectMysql()

	service.GlobalInit()

	fmt.Println(redis.Config)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	eee := e.Group("/blog/api")
	eee.GET("/article", controller.Article.List)

	eee.GET("/setting/:key", controller.Index.GetSetting)

	eee.GET("/article/:id", controller.Article.Info)
	eee.GET("/cate", controller.Cate.List)
	eee.GET("/GetFeatureSetting", controller.Index.GetFeatureSetting)

	e.Logger.Fatal(e.Start(":1323"))
}
