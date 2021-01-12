package user

import (
	"GoBlog/admin/controller/response"
	"GoBlog/admin/controller/service"
	"GoBlog/lib/helper"
	"GoBlog/logic/model"
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Userinfo struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

type UserDetail struct {
	model.SystemUser
	CheckedRoles []string `json:"checkedRoles"`
}

type userController struct{}

var User = &userController{}

func (*userController) Info(c echo.Context) error {
	session := c.Get("session").(service.Session)
	user, _ := session.GetSystemUser()
	var results []string
	model.SystemUserRole{}.Model().Select("system_role.name").Joins("LEFT join system_role  on system_user_role.system_role_id = system_role.id ").Where("system_user_role.system_user_id = ?", user.Id).Where("system_role.status = 1").Pluck("system_role.name", &results)
	var info Userinfo
	info.Roles = results
	info.Name = user.Name
	info.Avatar = user.Avatar
	info.Introduction = user.Introduction
	return c.JSON(http.StatusOK, response.ShowData(info))
}

func (*userController) Index(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	var results []model.SystemUser
	count := 0
	model.SystemUser{}.Model().Where("status = 1").Count(&count).Offset((page - 1) * limit).Limit(limit).Order(" id desc").Scan(&results)

	data := make(map[string]interface{})
	data["items"] = results
	data["total"] = count
	return c.JSON(http.StatusOK, response.ShowData(data))
}

func (*userController) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("id"))
	}
	user := model.SystemUser{}
	model.SystemUser{}.Model().Model(&user).Where("id = ?", id).Update("status", "0")
	return c.JSON(http.StatusOK, response.ShowData("success"))
}

func (*userController) Search(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("name"))
	}
	user := model.SystemUser{}
	res := user.GetAllByName(name)
	nameList := make(map[string][]map[string]interface{}, 0)
	nameList["items"] = res
	return c.JSON(http.StatusOK, response.ShowData(nameList))

}

func (*userController) Detail(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("id"))
	}

	user := model.SystemUser{}
	model.SystemUser{}.Model().Where("id = ?", id).First(&user)

	if user.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrormsg("user_error"))
	}

	var role []string
	model.SystemUserRole{}.Model().Select("system_role.name").Joins("LEFT join system_role  on system_user_role.system_role_id = system_role.id ").Where("system_user_role.system_user_id = ?", user.Id).Where("system_role.status = 1").Pluck("system_role.name", &role)

	detail := UserDetail{}
	detail.CheckedRoles = role
	detail.Id = user.Id
	detail.Name = user.Name
	detail.Nickname = user.Nickname
	detail.Phone = user.Phone
	detail.Status = user.Status

	return c.JSON(http.StatusOK, response.ShowData(detail))
}

func (*userController) Create(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if _, ok := data["name"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	if _, ok := data["nickname"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	if _, ok := data["password"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	if _, ok := data["repassword"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["status"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	userModel := model.SystemUser{}
	model.SystemUser{}.Model().Where("name = ? ", data["name"].(string)).First(&userModel)
	userModel.Name = data["name"].(string)
	if userModel.Id != 0 {
		return c.JSON(http.StatusOK, response.ShowErrormsg("name_exists"))
	}
	userModel.Name = data["name"].(string)
	if data["password"].(string) != data["repassword"].(string) {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	userModel.Password = data["password"].(string)
	userModel.Salt = helper.GetRandomBoth(4)
	userModel.Password = helper.Sha1En(userModel.Password + userModel.Salt)
	userModel.Name = data["name"].(string)
	userModel.Nickname = data["nickname"].(string)
	if _, ok := data["phone"]; ok {
		userModel.Phone = data["phone"].(string)
	}
	if _, ok := data["status"]; ok && data["status"].(bool) {
		userModel.Status = 1
	}
	userModel.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	userModel.Ctime = helper.JSONTime{}.Create()
	if _, ok := data["checkedRoles"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	roles := data["checkedRoles"].([]interface{})
	id, _ := userModel.Add(roles)
	if id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	return c.JSON(http.StatusOK, response.ShowData(userModel))
}

func (*userController) Repasswd(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if _, ok := data["id"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	userModel := model.SystemUser{}
	model.SystemUser{}.Model().Where("id = ?", data["id"]).First(&userModel)

	if userModel.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrormsg("user_error"))
	}

	if userModel.Name == "admin" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["password"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["repassword"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	userModel.Password = data["password"].(string)
	if userModel.Password != data["repassword"].(string) {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	userModel.Salt = helper.GetRandomBoth(4)
	userModel.Password = helper.Sha1En(userModel.Password + userModel.Salt)

	model.SystemUser{}.Model().Model(&userModel).Updates(userModel)
	return c.JSON(http.StatusOK, response.ShowData(userModel))
}

func (*userController) Edit(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["id"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	userModel := model.SystemUser{}
	model.SystemUser{}.Model().Where("id = ?", data["id"]).First(&userModel)

	if userModel.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrormsg("user_error"))
	}

	if _, ok := data["nickname"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	if _, ok := data["status"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	if _, ok := data["status"]; ok && data["status"].(bool) {
		userModel.Status = 1
	} else {
		userModel.Status = 0
	}
	userModel.Nickname = data["nickname"].(string)
	if _, ok := data["phone"]; ok {
		userModel.Phone = data["phone"].(string)
	}
	if _, ok := data["checkedRoles"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	roles := data["checkedRoles"].([]interface{})

	if !userModel.Update(roles) {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))

	}
	return c.JSON(http.StatusOK, response.ShowData(userModel))
}
