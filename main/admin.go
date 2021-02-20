package main

import (
	"GoBlog/admin/controller"
	"GoBlog/admin/controller/menu"
	"GoBlog/admin/controller/middleware"
	"GoBlog/admin/controller/role"
	"GoBlog/admin/controller/user"
	_ "GoBlog/config"
	"GoBlog/lib/es"
	_ "GoBlog/lib/helper"
	"GoBlog/lib/mysql"
	_ "encoding/json"
	_ "fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	eeem "github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	/*
		mysql.ConnectMysql(config.MySQL, "default")
		defer mysql.DisconnectMysql()
	*/
	es.InitEsConnect()
	mysql.InitMysqlConnect()

	e := echo.New()
	e.Use(eeem.Static("/upload"))

	e.POST("/admin/upload/image", controller.Index.ImgUpload)

	e.GET("/admin/del/image", controller.Index.DelImage)

	e.GET("/admin/showimage", controller.Index.ReadImage)

	//aaaaaaa := helper.JsonSqlValue{}.Create(map[string]interface{}{"dsklvmds": "mbfkgbmfglkbmfl", "vsvfdvfd": 394930})
	///fmt.Println(aaaaaaa)

	//aaaa := echomiddleware.CORS()

	//	e.Use(aaaa)

	admin := e.Group("/admin", middleware.CheckToken)

	e.Validator = &CustomValidator{validator: validator.New()}
	//admin.GET("/login", user.Login.Index)

	e.GET("/aaaaa", controller.Index.GetIndexPics)

	e.GET("/admin/login", user.Login.Index)
	e.POST("/admin/login", user.Login.Index)

	admin.POST("/logout", user.Login.Logout)
	admin.GET("/info", user.User.Info)
	admin.GET("/routes", menu.Menu.List)
	admin.GET("/menu", menu.Menu.Index)
	admin.POST("/menu", menu.Menu.Create)
	admin.PUT("/menu", menu.Menu.Edit)

	admin.DELETE("/menu", menu.Menu.Delete)
	admin.GET("/role/list", menu.Menu.Roles)

	admin.GET("/user", user.User.Index)
	admin.GET("/user/delete", user.User.Delete)

	admin.GET("/dashboard", menu.Menu.Dashboard)

	admin.GET("/user/search", user.User.Search)

	admin.GET("/user/detail", user.User.Detail)

	admin.POST("/user/repasswd", user.User.Repasswd)

	admin.POST("/user/create", user.User.Create)
	admin.POST("/user/edit", user.User.Edit)

	admin.GET("/role/index", role.Role.Index)
	admin.POST("/role/add", role.Role.AddRole)
	admin.POST("/role/delete", role.Role.DeleteRole)
	admin.POST("/role/update", role.Role.UpdateRole)

	admin.GET("/article/list", controller.Article.List)
	admin.GET("/article/detail", controller.Article.Info)
	admin.POST("/article/update", controller.Article.Update)
	admin.POST("/article/create", controller.Article.Create)

	admin.GET("/cate/info", controller.Cate.Info)
	admin.GET("/cate/list", controller.Cate.List)
	admin.POST("/cate/create", controller.Cate.Create)
	admin.POST("/cate/update", controller.Cate.Update)

	admin.POST("/UpdateIndexPics", controller.Index.UpdateIndexPics)
	admin.GET("/GetIndexPics", controller.Index.GetIndexPics)

	admin.POST("/PostWebsiteInfo", controller.Index.PostWebsiteInfo)

	admin.GET("/GetWebsiteInfo", controller.Index.GetWebsiteInfo)

	e.GET("/admin/pong", func(c echo.Context) error {
		return c.JSON(200, "{'message':'pong'}")
	})

	e.Logger.Fatal(e.Start(":888"))
}
