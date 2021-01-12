package menu

import (
	"GoBlog/admin/controller/response"
	"GoBlog/admin/controller/service"
	"GoBlog/lib/helper"
	"GoBlog/logic/model"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	_ "time"
)

type Role struct {
	Id          uint32        `form:"id" json:"id"`
	Name        string        `form:"name" json:"name"`
	Description string        `form:"description" json:"description"`
	Status      int8          `form:"status" json:"status"`
	Routes      []interface{} `form:"routes" json:"routes"`
}

type menuController struct{}

var Menu = &menuController{}

func (*menuController) List(c echo.Context) error {
	var menuArr []model.SystemMenu
	model.SystemMenu{}.Model().Find(&menuArr)
	var menuMap = make(map[uint32][]model.SystemMenu, 0)
	role := model.SystemRole{}
	mrArr := role.GetRowMenu()
	for _, value := range menuArr {
		value.Hidden = 0
		menuMap[value.Pid] = append(menuMap[value.Pid], value)
	}
	jsonArr := model.SystemMenu{}.TreeMenuNew(menuMap, 0, mrArr)
	return c.JSON(http.StatusOK, response.ShowData(jsonArr))
}

func (*menuController) Index(c echo.Context) error {
	var menuArr []model.SystemMenu
	model.SystemMenu{}.Model().Find(&menuArr)

	var menuMap = make(map[uint32][]model.SystemMenu, 0)
	for _, value := range menuArr {
		menuMap[value.Pid] = append(menuMap[value.Pid], value)
	}
	var menuNewArr []model.SystemMenu
	menuNewArr = model.SystemMenu{}.TreeNode(menuMap, 0)
	return c.JSON(http.StatusOK, response.ShowData(menuNewArr))
}

func (*menuController) Dashboard(c echo.Context) error {
	session := c.Get("session").(service.Session)
	user, _ := session.GetSystemUser()
	var menuArr []model.SystemMenu
	if user.Name == "admin" {
		model.SystemMenu{}.Model().Where("status = 1").Find(&menuArr)
	} else {
		model.SystemMenu{}.Model().Select(" distinct  system_menu.*").
			Joins("join system_role_menu system_role_menu.system_menu_id=  system_menu.id").
			Joins("join system_user_role  system_role_menu.system_role_id=  system_user_role.system_role_id").
			Where("system_menu.status = ?", 1).
			Where("system_user_role.system_user_id = ?", user.Id).
			Find(menuArr)
	}
	var menuMap = make(map[uint32][]model.SystemMenu, 0)
	for _, value := range menuArr {
		menuMap[value.Pid] = append(menuMap[value.Pid], value)
	}

	role := model.SystemRole{}
	mrArr := role.GetRowMenu()
	jsonStr := model.SystemMenu{}.TreeMenuNew(menuMap, 0, mrArr)
	return c.JSON(http.StatusOK, response.ShowData(jsonStr))
}

func (*menuController) Roles(c echo.Context) error {
	menu := model.SystemMenu{}
	var roleArr []model.SystemRole
	model.SystemRole{}.Model().Find(&roleArr)
	var roleMenu []Role
	for _, value := range roleArr {
		r := Role{}
		r.Id = value.Id
		r.Name = value.Name
		r.Status = value.Status
		r.Description = value.Description
		menuArr := menu.GetRouteByRole(value.Id)
		if menuArr != nil {
			var menuMap = make(map[uint32][]model.SystemMenu, 0)
			for _, value := range menuArr {
				value.Hidden = 0
				menuMap[value.Pid] = append(menuMap[value.Pid], value)
			}
			role := model.SystemRole{}
			mrArr := role.GetRowMenu()
			jsonStr := model.SystemMenu{}.TreeMenuNew(menuMap, 0, mrArr)
			if jsonStr != nil {
				r.Routes = jsonStr
			}
		}
		roleMenu = append(roleMenu, r)
	}
	return c.JSON(http.StatusOK, response.ShowData(roleMenu))
}

func (*menuController) Create(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["name"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["path"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["component"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if _, ok := data["url"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	menus := model.SystemMenu{}
	model.SystemMenu{}.Model().Where("path = ?", data["path"].(string)).First(&menus)
	if menus.Id != 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("path不可重复"))
	}

	menus.Path = data["path"].(string)
	if menus.Path == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	menus.Name = data["name"].(string)

	menus.MetaTitle = data["name"].(string)
	menus.Component = data["component"].(string)
	menus.Url = data["url"].(string)
	menus.Redirect = data["redirect"].(string)
	menus.MetaIcon = data["meta_icon"].(string)
	if data["alwaysshow"].(bool) {
		menus.Alwaysshow = 1
	}
	if data["hidden"].(bool) {
		menus.Hidden = 1
	}
	if data["status"].(bool) {
		menus.Status = 1
	}

	menus.Ctime = helper.JSONTime{}.Create()

	if _, ok := data["sort"]; ok {
		menus.Sort, _ = strconv.Atoi(helper.Strval(data["sort"]))
	}

	if _, ok := data["pid"]; ok {
		pid, _ := strconv.Atoi(helper.Strval(data["pid"]))
		if pid != 0 {
			pmmmmm := model.SystemMenu{}
			model.SystemMenu{}.Model().Where("id = ?", pid).Where("status = 1").Where("hidden = 0").First(&pmmmmm)
			if pmmmmm.Id == 0 {
				return c.JSON(http.StatusOK, response.ShowErrorParams("pid错误"))
			}
			menus.Pid = uint32(pid)
		}
	}

	if menus.Pid == 0 {
		menus.Level = 0
	} else {
		pidMenuModel := model.SystemMenu{}
		model.SystemMenu{}.Model().Where("id = ?", menus.Pid).First(&pidMenuModel)
		menus.Level = pidMenuModel.Level + 1
	}
	menus.Model().Create(&menus)
	if menus.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	return c.JSON(http.StatusOK, response.ShowData(menus))
}

func (*menuController) Edit(c echo.Context) error {

	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	if _, ok := data["id"]; !ok {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	menu := model.SystemMenu{}

	model.SystemMenu{}.Model().Where("id=?", data["id"]).First(&menu)

	if menu.Id == 0 {
		return c.JSON(http.StatusOK, response.ShowErrorParams("要修改数据不存在"))
	}

	if _, ok := data["name"]; ok {
		menu.Name = data["name"].(string)
		menu.MetaTitle = data["name"].(string)

	}
	if _, ok := data["path"]; ok {
		menu.Path = data["path"].(string)
	}
	if _, ok := data["component"]; ok {
		menu.Component = data["component"].(string)

	}
	if _, ok := data["url"]; ok {
		menu.Url = data["url"].(string)
	}

	if data["path"].(string) == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}

	menus := model.SystemMenu{}
	model.SystemMenu{}.Model().Where("path = ?", data["path"].(string)).First(&menus)
	if menus.Id != 0 && menus.Id != menu.Id {
		return c.JSON(http.StatusOK, response.ShowErrorParams("path不可重复"))
	}

	menu.Redirect = data["redirect"].(string)
	menu.MetaIcon = data["meta_icon"].(string)
	if data["meta_nocache"].(bool) {
		menu.MetaNocache = 1
	}
	if data["hidden"].(bool) {
		menu.Hidden = 1
	}
	if data["alwaysshow"].(bool) {
		menu.Alwaysshow = 1
	}
	if data["status"].(bool) {
		menu.Status = 1
	}

	if _, ok := data["sort"]; ok {
		menu.Sort, _ = strconv.Atoi(helper.Strval(data["sort"]))
	}

	if _, ok := data["pid"]; ok {
		pid, _ := strconv.Atoi(helper.Strval(data["pid"]))
		if pid != 0 {
			pmmmmm := model.SystemMenu{}
			model.SystemMenu{}.Model().Where("id = ?", pid).Where("status = 1").Where("hidden = 0").First(&pmmmmm)
			if pmmmmm.Id == 0 {
				return c.JSON(http.StatusOK, response.ShowErrorParams("pid错误"))
			}
			menu.Pid = uint32(pid)
		}
	}

	model.SystemMenu{}.Model().Model(&menu).Updates(menu)
	return c.JSON(http.StatusOK, response.ShowData("success"))
}

func (*menuController) Delete(c echo.Context) error {
	str, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	id := string(str)
	if id == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	mm := model.SystemMenu{}
	dddd, _ := strconv.Atoi(id)
	mm.Id = uint32(dddd)
	model.SystemMenu{}.Model().Delete(&mm)
	return c.JSON(http.StatusOK, response.ShowData("success"))
}
