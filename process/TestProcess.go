package process

import (
     "fmt"
      "time"
      "GoBlog/lib/amqp"
      "GoBlog/lib/websocket"

)


type TestProcess struct {
	 
}

func (l *TestProcess) Run(){
	 //	  fmt.Printf("士大夫反对实力派【1jsadoifjdsdskfjdslkfjdskljfdskl")
    conn :=	amqp.GetAmqp() 
    if conn != nil {
    	 defer conn.Close()
    }


websocket.WebsocketBase.ServerInit()
uuid:=websocket.WebsocketBase.OpenInit()
kk := make(map[string]interface{})
kk["dsadsa"] = "dsadsads"
websocket.WebsocketBase.SetSession(uuid,kk)

kk["erwrwrwerwsa"] = "dsadsads"

websocket.WebsocketBase.SetSession(uuid,kk)

websocket.WebsocketBase. UpdateUUidDataTime(uuid  )

fmt.Println( "uuid:  ", websocket.WebsocketBase.GetUUidData(uuid))
websocket.WebsocketBase.SendTo(uuid,"u89etu894ut4389")

	websocket.WebsocketBase.Disconnect(uuid)
 
  //  websocket.WebsocketBase. ServerInit()
	  	  time.Sleep(3* time.Second)
	  	  // panic("test panic")
}
