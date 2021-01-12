package config

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

var MySQL = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "127.0.0.1", 3306, "test")

var MysqlConfig = map[string]map[string]string{}

func init() {
	content, err := ioutil.ReadFile("env/mysql.json")
	if err != nil {
		//log.Fatal(err)
		return
	}
	configStr := string(content)
	gjson.Parse(configStr).ForEach(func(key, value gjson.Result) bool {
		vvv := map[string]string{}
		value.ForEach(func(k, v gjson.Result) bool {
			vvv[k.String()] = v.String()
			return true
		})
		vvv["conString"] = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", vvv["user"], vvv["password"], vvv["host"], vvv["port"], vvv["db"])
		MysqlConfig[key.String()] = vvv
		return true
	})

}
