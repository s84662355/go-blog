package controller

import (
	"GoBlog/admin/controller/response"
	_ "GoBlog/admin/controller/service"
	"GoBlog/lib/helper"
	"GoBlog/logic/model"
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	_ "strconv"
	"strings"
	"time"
)

const indexPics = "IndexPics"
const websiteInfo = "WebsiteInfo"

type indexController struct{}

var Index = &indexController{}

func (*indexController) GetWebsiteInfo(c echo.Context) error {
	setting := model.Setting{}
	setting.Model().Where("key_name = ?", websiteInfo).First(&setting)
	return c.JSON(http.StatusOK, response.ShowData(setting.Content))
}

func (*indexController) PostWebsiteInfo(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if v, ok := data["websiteInfo"]; ok {
		setting := model.Setting{}
		setting.KeyName = websiteInfo
		setting.Model().Table(setting.TableName()).Where("key_name = ?", websiteInfo).First(&setting)
		setting.Content = helper.JsonSqlValue{}.Create(v)
		if setting.Id != 0 {
			setting.Model().Table(setting.TableName()).Save(&setting)
		} else {
			setting.Model().Table(setting.TableName()).Create(&setting)
		}
		return c.JSON(http.StatusOK, response.ShowData(setting.Content))
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams("缺少参数"))
	}
}

func (*indexController) GetIndexPics(c echo.Context) error {
	setting := model.Setting{}
	setting.Model().Where("key_name = ?", indexPics).First(&setting)
	return c.JSON(http.StatusOK, response.ShowData(setting.Content))
}

func (*indexController) UpdateIndexPics(c echo.Context) error {
	jsonstr, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonstr, &data)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	if v, ok := data["imgs"]; ok {
		imgs := v.([]interface{})
		setting := model.Setting{}
		setting.KeyName = indexPics
		setting.Model().Table(setting.TableName()).Where("key_name = ?", indexPics).First(&setting)
		setting.Content = helper.JsonSqlValue{}.Create(imgs)
		if setting.Id != 0 {
			setting.Model().Table(setting.TableName()).Save(&setting)
		} else {
			setting.Model().Table(setting.TableName()).Create(&setting)
		}
		return c.JSON(http.StatusOK, response.ShowData(setting.Content))
	} else {
		return c.JSON(http.StatusOK, response.ShowErrorParams("缺少参数"))
	}
}

//上传图片
func (*indexController) ImgUpload(c echo.Context) error {
	//获取表单数据 参数为name值
	f, err := c.FormFile("file")
	//错误处理
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
	}
	//根据当前时间生成目录
	timestamp := time.Now().Unix()
	tm2 := time.Unix(timestamp, 0)
	//创建上传目录
	uploadDir := "upload/admin/"
	relative := tm2.Format("20060102") + "/"
	os.MkdirAll(uploadDir+relative, os.ModePerm)
	//扩展名
	ext := path.Ext(f.Filename)
	relative = f.Filename
	src, err := f.Open()
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("文件失败"))
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(uploadDir + relative)
	if err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("创建文件失败"))
	}
	defer func() {
		dst.Close()
		////go os.Remove(uploadDir + relative)
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusOK, response.ShowErrorParams("转移文件失败"))
	}
	//aaaa := helper.AliOss{}.Create()
	//res := aaaa.Upload(uploadDir+relative, helper.GetRandomBoth(32)+ext)

	//保存成功返回正确的Json数据
	//data := "/admin/showimage?imgname=upload/admin/" + relative
	///data := "/admin/" + relative
	return c.JSON(http.StatusOK, response.ShowData(res))
}

func (*indexController) ReadImage(c echo.Context) error {
	//imgname=upload/20201225/YC43nmZPNYlrYdIv36bBDwBpA3DJBE3d.ico
	url := c.QueryParam("imgname")
	url = helper.SubstrContains(url, "upload/")
	//顾虑危险字符
	url = strings.Replace(url, "../", "", -1)
	f, _ := ioutil.ReadFile(url)
	return c.String(http.StatusOK, string(f))
}

//删除图片
func (*indexController) DelImage(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusOK, response.ShowErrorParams("url"))
	}
	url = helper.SubstrContains(url, "upload/")
	//顾虑危险字符
	url = strings.Replace(url, "../", "", -1)
	if helper.IsFile(url) {
		err := os.Remove(url)
		if err != nil {
			return c.JSON(http.StatusOK, response.ShowErrorParams("fail"))
		}
	}
	return c.JSON(http.StatusOK, response.ShowData("success"))
}
