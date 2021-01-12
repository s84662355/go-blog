package user

import (
    "GoBlog/admin/controller/response"
    "GoBlog/admin/controller/service"
    "GoBlog/lib/helper"
    _ "GoBlog/lib/redis"
    "GoBlog/logic/model"
    "fmt"
    "github.com/labstack/echo/v4"
    "net/http"
)

type UserRequest struct {
    Username string `json:"username" form:"username" query:"username" validate:"required"`
    Password string `json:"password" form:"password" query:"password" validate:"required"`
}

type loginController struct{}

var Login = &loginController{}

func (*loginController) Index(c echo.Context) error {
    fmt.Println("dfsfds")
    u := new(UserRequest)
    if err := c.Bind(u); err != nil {
        fmt.Println(err)
        return c.JSON(http.StatusBadRequest, response.ShowErrorParams("参数错误"))
    }
    if err := c.Validate(u); err != nil {
        fmt.Println(err)
        return c.JSON(http.StatusBadRequest, response.ShowErrorParams("参数错误"))
    }
    results := model.SystemUser{}
    model.SystemUser{}.Model().Where(" name  =  '" + u.Username + "'").Where("status = 1").First(&results)
    if results.Id == 0 {
        return c.JSON(http.StatusOK, response.ShowErrormsg("用户不存在"))
    }
    if helper.Sha1En(u.Password+results.Salt) != results.Password {
        return c.JSON(http.StatusOK, response.ShowErrormsg("密码错误"))
    }
    token := "admin_" + helper.GetRandomString(20)
    ss := service.CreateSession(token)
    if ss.SetSystemUser(results) {
        return c.JSON(http.StatusOK, response.ShowData(map[string]interface{}{"token": token}))
    }
    return c.JSON(http.StatusOK, response.ShowErrormsg("登陆失败"))
}

func (*loginController) Logout(c echo.Context) error {
    session := c.Get("session").(service.Session)
    session.LoginOut()
    return c.JSON(http.StatusOK, response.ShowSuccess("success"))
}
