package controller

import (
	"GoBlog/admin/controller/response"
	_ "GoBlog/admin/controller/service"
	_ "GoBlog/lib/helper"
	"GoBlog/logic/model"
	_ "encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	_ "io/ioutil"
	"net/http"
	_ "strconv"
)

type cateController struct{}

var Cate = &cateController{}

func (*cateController) Create(c echo.Context) error {
	mm := model.Cate{}

	if c.QueryParam("name") != "" {
		mm.Name = c.QueryParam("name")
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if !mm.Create() {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	return c.JSON(http.StatusOK, response.ShowData(mm))
}

func (*cateController) Update(c echo.Context) error {
	mm := model.Cate{}.Get(c.QueryParam("id"))
	if mm.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("分类不存在"))
	}

	if c.QueryParam("name") != "" {
		mm.Name = c.QueryParam("name")
	}
	mm.Update()
	return c.JSON(http.StatusOK, response.ShowData(mm))
}

func (*cateController) Info(c echo.Context) error {
	mm := model.Cate{}.Get(c.QueryParam("id"))
	return c.JSON(http.StatusOK, response.ShowData(mm))
}

func (*cateController) List(c echo.Context) error {
	res := model.Cate{}.All()
	return c.JSON(http.StatusOK, response.ShowData(res))
}

func (*cateController) Delete(c echo.Context) error {
	mm := model.Cate{}.Get(c.QueryParam("id"))
	if mm.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	count := 0
	model.Article{}.Model().Where("cate_id = ?", mm.Id).Where("deleted_at is null").Count(&count)

	if count != 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("删除失败"))
	}
	mm.Delete()
	return c.JSON(http.StatusOK, response.ShowData("成功"))
}
