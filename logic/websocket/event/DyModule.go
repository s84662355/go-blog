package event

import (
	"GoBlog/lib/websocket"
	"fmt"
)

type DyModule struct {
     
}

type DDD struct{
		A string   `jpath:"a"`

}


var dyModule *DyModule = new(DyModule)


func init(){
	 Register("Dy","user_info",dyModule.user_info)
}

 
 

func (l *DyModule) user_info(data map[string]interface{},session *websocket.WebsocketSession){
     fmt.Println("user_info")

ddd:=new(DDD)
GetParam(data,ddd)


      fmt.Println("集散地欧冠降低偶发",ddd)

       fmt.Println(session )
}