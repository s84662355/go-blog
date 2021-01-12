package role

import (
	"GoBlog/admin/controller/response"
	_ "GoBlog/admin/controller/service"
	"GoBlog/lib/helper"
	"GoBlog/logic/model"
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	_ "strconv"
)

type roleController struct{}

var Role = &roleController{}

func (*roleController) UpdateRole(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["id"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	mm := model.SystemRole{}
	id := uint32(data["id"].(float64))
	model.SystemRole{}.Model().Where("name = ?", data["name"]).First(&mm)
	if mm.Id != 0 && mm.Id != id {
		return c.JSON(http.StatusOK, response.ShowErrorParams("role_error"))
	}
	mm.Id = id
	mm.Name = data["name"].(string)
	mm.AliasName = data["name"].(string)
	mm.Description = data["description"].(string)
	if _, ok := data["status"]; ok {
		if data["status"].(bool) {
			mm.Status = 1
		}
	}
	var ids []uint32
	if _, ok := data["routes"]; ok {
		ids = model.SystemRole{}.TreeRoutes(data["routes"].([]interface{}))
	}
	if !mm.Update(ids) {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	datas := map[string]string{"status": "success"}
	return c.JSON(http.StatusOK, response.ShowData(datas))
}

func (*roleController) AddRole(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	mm := model.SystemRole{}

	model.SystemRole{}.Model().Where("name = ?", data["name"].(string)).First(&mm)
	if mm.Id != 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("role_error"))
	}
	mm.Name = data["name"].(string)
	mm.AliasName = data["name"].(string)
	mm.Description = data["description"].(string)
	if _, ok := data["status"]; ok {
		if data["status"].(bool) {
			mm.Status = 1
		}
	}
	mm.Ctime = helper.JSONTime{}.Create()

	if !mm.AddCommit(data["routes"].([]interface{})) {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	return c.JSON(http.StatusOK, response.ShowData(mm))
}

func (*roleController) DeleteRole(c echo.Context) error {

	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if _, ok := data["id"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	id := int(data["id"].(float64))
	role := model.SystemRole{}
	model.SystemRole{}.Model().Where("id = ?", id).First(&role)

	if role.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	role.Delete()
	if !role.Delete() {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	datas := map[string]string{"status": "success"}
	return c.JSON(http.StatusOK, response.ShowData(datas))
}

func (*roleController) Index(c echo.Context) error {
	roles := model.SystemRole{}
	list := roles.GetNameList()
	return c.JSON(http.StatusOK, response.ShowData(list))
}
