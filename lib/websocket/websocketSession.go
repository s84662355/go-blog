package websocket

import (
	"GoBlog/lib/redis"
 	goredis  "github.com/go-redis/redis"
	"strconv"
	"fmt"
	"time"
)

const SESSION_DEVICE = "session_device"

const DEVICE_ZADD = "session_device_zadd"

type WebsocketSession struct{
   uuid string
}

func NewWebsocketSession(uuid string) *WebsocketSession{
     v := new(WebsocketSession)
     v.uuid = uuid
     return v
}


/*
 *
 * 设备登录
 *
 * */
func (l *WebsocketSession) Login(device_id string, session_data map[string]interface{}){
	client := redis.GetRedis()
	t := time.Now()
	score , _ := strconv.ParseFloat(fmt.Sprintf("%d",t.Unix()), 64)

	client.ZAdd(DEVICE_ZADD, goredis.Z{Score:score,Member:device_id}).Result()

	client.Set(fmt.Sprintf("%s%s",SESSION_DEVICE , device_id),l.uuid,EXPIRE * time.Second)

	WebsocketBase.SetSession(l.uuid,session_data)

	WebsocketBase.UpdateUUidDataTime(l.uuid)

	/*
	SwoftRedis::connection(BaseData::SESSION_REDIS)->set(self::SESSION_DEVICE.$device_id,$this->uuid);
	SwoftRedis::connection(BaseData::SESSION_REDIS)->expire(self::SESSION_DEVICE.$device_id,BaseData::EXPIRE);
    */
}


func (l *WebsocketSession) LoginOut(){
    l.LoginOutByUuid(l.uuid)
}

func (l *WebsocketSession) LoginOutByUuid(uuid string){
	sessionData := WebsocketBase.GetSession(uuid)
	device_id , ok := sessionData["device_id"]
	if ok {
		client := redis.GetRedis()
		deviceIdString , ok := device_id.(string)
		if ok {
			client.Del(fmt.Sprintf("%s%s",SESSION_DEVICE,deviceIdString))
			WebsocketBase.Disconnect(uuid)
		}
	}
}

func (l *WebsocketSession) Close(){
	WebsocketBase.Disconnect(l.uuid)
}

func (l *WebsocketSession) Islogin(device_id string) string{
	client := redis.GetRedis()
    result , err := client.Get(fmt.Sprintf("%s%s",SESSION_DEVICE , device_id)).Result()
    if err != nil{
    	   return ""
	}
	return result
}

func (l *WebsocketSession) Conflict(device_id string) int64{
	client := redis.GetRedis()
	incr , err := client.Incr(fmt.Sprintf("%s%s",device_id,"__conflict")).Result()
	if err != nil{
		return 0
  	}

	client.Expire(fmt.Sprintf("%s%s",device_id,"__conflict"),15*time.Second)
	return  incr

	//$incr = SwoftRedis::connection(BaseData::SESSION_REDIS)->incr($device_id.'__conflict');
	//SwoftRedis::connection(BaseData::SESSION_REDIS)->expire($device_id.'__conflict',15)

}

func (l *WebsocketSession) CheckHreat() bool{
	uuid := l.uuid
	sessionData := WebsocketBase.GetSession(uuid)
	if len(sessionData) == 0 {
		return  false
	}

	device_id , ok := sessionData["device_id"]
	client := redis.GetRedis()
	var device_uuid string = ""

	if ok {
 		deviceIdString , ok := device_id.(string)
		if ok {
			device_uuid = client.Get(fmt.Sprintf("%s%s",SESSION_DEVICE,deviceIdString)).Val()

 		}else {
 			return false
		}

	}else {
		return false
	}

	if device_uuid != "" && device_uuid != uuid{
		return false
	}

	WebsocketBase.UpdateUUidDataTime(uuid)
	client.Expire(fmt.Sprintf("%s%s",SESSION_DEVICE,device_id),EXPIRE * time.Second)
	return true
}


func (l *WebsocketSession) IsEstablished() bool{
	device_id :=  l.GetDeviceId()
	client := redis.GetRedis()
	device_uuid , err := client.Get(fmt.Sprintf("%s%s",SESSION_DEVICE,device_id)).Result()
	if err != nil {
		return  false
	}
	if device_uuid != l.uuid{
		return false
	}
	return true
}

func (l *WebsocketSession) GetDeviceId() string{
	uuid := l.uuid
	sessionData := WebsocketBase.GetSession(uuid)
	if len(sessionData) == 0 {
		return  ""
	}

	device_id , ok := sessionData["device_id"]

	if ok {
		deviceIdString, ok := device_id.(string)
		if ok{
			return deviceIdString
		}
	}
	return ""
}

func (l *WebsocketSession) Session(sessionData map[string]interface{}) map[string]interface{}{
	WebsocketBase.SetSession(l.uuid,sessionData)
	return WebsocketBase.GetSession(l.uuid)
}


func (l *WebsocketSession) IsBind(isLoginOut bool) bool{
	device_id := l.GetDeviceId()
	if device_id != "" {
		if isLoginOut{
			l.LoginOut()
		}

		return false
	}else{
		return  true
	}
}

func (l *WebsocketSession) GetAppDeviceId() int{
	uuid := l.uuid
	sessionData := WebsocketBase.GetSession(uuid)
	id , ok := sessionData["id"]
	if ok {
		return  id.(int)
	}
	return 0
}




