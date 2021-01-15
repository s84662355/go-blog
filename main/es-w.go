package main

import (
	"GoBlog/es/model"
	"GoBlog/lib/es"
	"GoBlog/lib/mysql"
	_ "context"
	_ "fmt"
	"github.com/olivere/elastic"
)

var client *elastic.Client

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

//初始化
func init() {
	es.InitEsConnect()
	//errorlog := log.New(os.Stdout, "APP", log.LstdFlags)

	//这个地方有个小坑 不加上elastic.SetSniff(false) 会连接不上
	client = es.Es("default")

	mysql.InitMysqlConnect()

	//fmt.Printf("Elasticsearch version %s\n", esversion)
}

//创建
func create() {

	model.Article{}.PostDataById(5)

	///fmt.Println(aaa)

}

func main() {
	create()
}
