package controller

import (
	"GoBlog/admin/controller/response"
	_ "GoBlog/admin/controller/service"

	esm "GoBlog/es/model"
	"GoBlog/lib/helper"
	"GoBlog/logic/model"
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	_ "time"
)

type articleController struct{}

var Article = &articleController{}

func (*articleController) List(c echo.Context) error {
	params := map[string]interface{}{}
	if c.QueryParam("title") != "" {
		params["title"] = c.QueryParam("title")
	}

	if c.QueryParam("cate_id") != "" {
		params["cate_id"] = c.QueryParam("cate_id")
	}

	if c.QueryParam("status") != "" {
		params["status"] = c.QueryParam("status")
	}

	QueryParams := c.QueryParams()

	if v, ok := QueryParams["dateValue[]"]; ok {
		dateValue := v
		params["start_date"] = dateValue[0]
		params["end_date"] = dateValue[1]
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("limit"))

	esm.Article{}.Search(params, page, size)

	//	return c.JSON(http.StatusOK, response.ShowData(model.Article{}.List(params, page, size)))

	return c.JSON(http.StatusOK, response.ShowData(esm.Article{}.Search(params, page, size)))

}

func (*articleController) Info(c echo.Context) error {
	results := model.Article{}.Info(c.QueryParam("id"))
	return c.JSON(http.StatusOK, response.ShowData(results))
}

func (*articleController) Delete(c echo.Context) error {
	mm := model.Article{}.Get(c.QueryParam("id"))
	if mm.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	go func() {
		mm.Delete()
	}()

	return c.JSON(http.StatusOK, response.ShowData("成功"))
}

func (*articleController) Create(c echo.Context) error {
	mm := model.Article{}
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if v, ok := data["title"]; ok {
		mm.Title = v.(string)
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams(" title fail"))
	}

	if v, ok := data["image"]; ok {
		mm.Image = v.(string)
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams(" image fail"))
	}

	if v, ok := data["content"]; ok {
		mm.Content = helper.JsonSqlValue{}.Create(v)
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams(" content fail"))
	}

	if v, ok := data["summary"]; ok {
		mm.Summary = v.(string)
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams(" summary fail"))
	}

	if v, ok := data["cate_id"]; ok {
		cate := model.Cate{}.Get(v)
		if cate.Id == 0 {
			return c.JSON(http.StatusOK, response.ShowErrorParams("cate_id fail"))
		}

		switch v.(type) {
		case int:
			mm.CateId = v.(uint32)
			break
		case float64:
			mm.CateId = uint32(v.(float64))
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.CateId = uint32(vv)
			break
			break
		}
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams(" cate_id fail"))
	}
	if v, ok := data["sort"]; ok {
		switch v.(type) {
		case int:
			mm.Sort = v.(uint32)
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.Sort = uint32(vv)
			break
		}

	}
	if v, ok := data["status"]; ok {
		switch v.(type) {
		case int:
			mm.Status = v.(uint32)
			break
		case float64:
			mm.Status = uint32(v.(float64))
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.Status = uint32(vv)
			break
			break
		}
		if mm.Status != 0 && mm.Status != 1 {
			return c.JSON(http.StatusOK, response.ShowErrorParams("status fail"))
		}
	}

	if !mm.Create() {
		return c.JSON(http.StatusOK, response.ShowErrorParams("create fail"))
	}

	return c.JSON(http.StatusOK, response.ShowData(mm))
}

func (*articleController) Update(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	mm := model.Article{}.Get(data["id"])
	if mm.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if v, ok := data["image"]; ok {
		mm.Image = v.(string)
	}

	if v, ok := data["title"]; ok {
		mm.Title = v.(string)
	}
	if v, ok := data["content"]; ok {
		mm.Content = helper.JsonSqlValue{}.Create(v)
	}

	if v, ok := data["summary"]; ok {
		mm.Summary = v.(string)
	}

	if v, ok := data["cate_id"]; ok {
		cate := model.Cate{}.Get(v)
		if cate.Id == 0 {
			return c.JSON(http.StatusOK, response.ShowErrorParams("cate_id fail"))
		}
		switch v.(type) {
		case int:
			mm.CateId = v.(uint32)
			break
		case float64:
			mm.CateId = uint32(v.(float64))
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.CateId = uint32(vv)
			break
			break
		}
	}
	if v, ok := data["sort"]; ok {
		switch v.(type) {
		case int:
			mm.Sort = v.(uint32)
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.Sort = uint32(vv)
			break
		}
	}
	if v, ok := data["status"]; ok {
		switch v.(type) {
		case int:
			mm.Status = v.(uint32)
			break
		case float64:
			mm.Status = uint32(v.(float64))
			break
		case string:
			vv, _ := strconv.Atoi(v.(string))
			mm.Status = uint32(vv)
			break
			break
		}
		if mm.Status != 0 && mm.Status != 1 {
			return c.JSON(http.StatusOK, response.ShowErrorParams("status fail"))
		}
	}

	mm.Update()
	return c.JSON(http.StatusOK, response.ShowData(mm))
}
