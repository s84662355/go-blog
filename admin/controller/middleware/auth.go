package middleware

import (
	"GoBlog/admin/controller/response"
	"GoBlog/admin/controller/service"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

//中间件函数
func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-Token")
		if token == "" {
			token = c.QueryParam("token")
		}
		ss := service.CreateSession(token)
		_, err := ss.GetSystemUser()
		if err != nil {
			return c.JSON(http.StatusOK, response.ShowErrormsg("token错误"))
		}
		c.Set("session", ss)
		return next(c)
	}
}

func Cors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Response().Header().Set("Access-Control-Allow-Origin", "https://admin.duiniya.com, http://localhost:9529, http://localhost:9528, http://localhost:9527, http://localhost")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,OPTIONS,DELETE,UPDATE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "X-Requested-With,apitoken ,Origin, Content-Type, Cookie, X-CSRF-TOKEN, Accept, Authorization, X-XSRF-TOKEN,X-Token")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request().Method == "OPTIONS" {
			return c.NoContent(http.StatusNoContent)
		}

		return next(c)
	}
}
