package controller

import (
	_ "GoBlog/lib/helper"
	"GoBlog/logic/controller/response"
	"GoBlog/logic/model"
	"GoBlog/service"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ArticleRequest struct {
	Page int `json:"page" form:"page" query:"page" validate:"number,min=1"`
	//Size  int `json:"psize" form:"psize" query:"psize" validate:"number,min=10,max=20"`
	Title  *string `json:"title" form:"title" query:"title"`
	CateId int     `json:"cate_id" form:"cate_id" query:"cate_id"`
}

type ArticleController struct{}

var Article = &ArticleController{}

func (*ArticleController) List(c echo.Context) error {

	params := map[string]interface{}{}
	if c.QueryParam("title") != "" {
		params["title"] = c.QueryParam("title")
	}

	if c.QueryParam("cate_id") != "" {
		params["cate_id"] = c.QueryParam("cate_id")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	size := 10
	return c.JSON(http.StatusOK, response.Success(model.Article{}.List(params, page, size)))
}

func (*ArticleController) Info(c echo.Context) error {
	results := model.Article{}.Info(c.Param("id"))
	results.SetReadAmount()
	service.ReadCountQueue.Push(c.Param("id"), -1)
	return c.JSON(http.StatusOK, response.Success(results))
}
