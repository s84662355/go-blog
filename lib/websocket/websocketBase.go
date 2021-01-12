package websocket

import (
    "GoBlog/lib/helper"
    "GoBlog/config"
    "fmt"
    "time"
    "math/rand"
    "GoBlog/lib/redis"
    "encoding/json"
    goredis  "github.com/go-redis/redis"
    "strconv"
)

const WEBSOCKET_NODE_QUEUE_HSET = "websocket_node_queue_hset" ///每一个节点基本数据，例如redis的监听频道
const WEBSOCKET_CLIENT_HSET = "websocket_client_hset" //每一个连接的具体信息 key是uuid 每一个连接的唯一标识符
const WEBSOCKET_NODE_FD_HSET = "websocket_node_fd_hset" //当前fd对应的uuid
const DWSSERVER = "dwsserver"
const SESSION_REDIS = "WEBSOCKET_SESSION_REDIS" //用储存session信息的redis bean的key
const REDIS_CHAN = "WEBSOCKET_REDIS_CHAN"//分布式节点之间通讯的redis bean 的key
const WEBSOCKET_CLIENT_ZADD = "websocket_client_zadd"//存放所有uuid的集合
const OPENSERVER =  "OpenServer"
const MSGSERVER =  "MsgServer"
const CLOSESERVER = "CloseServer"
const PROCESSSERVER = "ProcessServer"
const SOCKETIOSERVER = "SocketIOServer"
const EXPIRE = 80
const TIMEOUT = 120

var WebsocketBase  = new(websocketBase)

type websocketBase struct{

}

func (l * websocketBase) GetNodeId() string{
    ip := helper.GetLocalIP()
    key :=   fmt.Sprintf("%s:%d",ip,config.WebSocketPort) 
    return  key 
}

func (l * websocketBase) CreateUUid() string{
    return   fmt.Sprintf("%s%d%d",l.GetNodeId(),time.Now().UnixNano(),rand.Intn(100000)) 
}

func (l * websocketBase) ServerInit(){
    key    := l.GetNodeId()
    queue  := fmt.Sprintf("%s%s",key,"_queue")
    client := redis.GetRedis() 
    hsetMap := make(map[string]interface{})
    hsetMap["queue"] = queue
    hsetMap["ip"] = helper.GetLocalIP()
    hsetMap["port"] = config.WebSocketPort
    hsetMapStringByte , _:=  json.Marshal(hsetMap)
    hsetMapString := string(hsetMapStringByte)
    fmt.Printf(hsetMapString)
    client.HSet(WEBSOCKET_NODE_QUEUE_HSET,key,hsetMapString)
}

func (l * websocketBase) OpenInit() string{
    client := redis.GetRedis()     
    uuid := l.CreateUUid()
	t := time.Now() 
    score , _ := strconv.ParseFloat(fmt.Sprintf("%d",t.Unix()), 64)

	client.ZAdd(WEBSOCKET_CLIENT_ZADD, goredis.Z{Score:score,Member:uuid}).Result()
   
    createDate := fmt.Sprintf("%d-%d-%d %d:%d:%d\n", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second())

    hmsetMap := make(map[string]interface{})
    hmsetMap["node"] = l.GetNodeId()
    hmsetMap["createDate"] = createDate
    hmsetMap["updateDate"] = createDate
    hmsetMap["updatetime"] = t.Unix()
    hmsetMap["session_data"] = "{}"
   
    client.HMSet(fmt.Sprintf( "%s%s" ,WEBSOCKET_CLIENT_HSET , uuid), hmsetMap).Result()
   
    client.Expire(fmt.Sprintf( "%s%s" ,WEBSOCKET_CLIENT_HSET , uuid), EXPIRE*time.Second) 
    return uuid
}


func (l * websocketBase) GetUUidData(uuid string) map[string]string{
    client := redis.GetRedis()    
    return client.HGetAll(fmt.Sprintf( "%s%s" ,WEBSOCKET_CLIENT_HSET , uuid)).Val();
}

func (l * websocketBase) GetNodeData(node string) map[string]interface{}{
	client := redis.GetRedis() 
    hsetMap := make(map[string]interface{})
    data := client.HGet(WEBSOCKET_NODE_QUEUE_HSET,node).Val()
    json.Unmarshal([]byte(data),&hsetMap)
    return  hsetMap 
}

func (l * websocketBase) GetSession(uuid string) map[string]interface{}{
    uuid_data := l.GetUUidData(uuid)

    hsetMap := make(map[string]interface{})
    if _ , ok := uuid_data["session_data"] ; ok{
          json.Unmarshal([]byte(uuid_data["session_data"]),&hsetMap)
    }

    return hsetMap
}

 
func (l * websocketBase) SetSession(uuid string, dataArr map[string]interface{}) bool{
    uuid_data := l.GetUUidData(uuid)
    if len(uuid_data)>0 {
		client := redis.GetRedis()     
		t := time.Now() 
	    createDate := fmt.Sprintf("%d-%d-%d %d:%d:%d\n", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second())
	    uuid_data["updateDate"] = createDate

	    session := l.GetSession(uuid)

	    for key   := range dataArr{
              session[key] = dataArr[key]
	    }

	    hashString , _ := json.Marshal(session)
	    uuid_data["session_data"] = string(hashString)
	    hashMap := make(map[string]interface{})

	    for k   := range uuid_data {
            hashMap[k] = uuid_data[k]
	    }
	    hashMap["updatetime"] =  t.Unix() 
	    client.HMSet(fmt.Sprintf("%s%s" ,WEBSOCKET_CLIENT_HSET , uuid), hashMap).Result()
		client.Expire(fmt.Sprintf( "%s%s" ,WEBSOCKET_CLIENT_HSET , uuid), EXPIRE*time.Second)
		return true
    }
    return false
}


func (l * websocketBase) UuidIsMatch(uuid string) bool{
        uuid_data :=  l.GetUUidData(uuid)
        if len(uuid_data) == 0{
               return false
        }
           
        if uuid_data["node"] != l.GetNodeId(){
            return false
        }
        return true
}


func (l * websocketBase) UpdateUUidDataTime(uuid string) bool{
       return l.SetSession(uuid,make(map[string]interface{}))
}


func (l * websocketBase) SendTo(uuid string, data string) bool{
    uuid_data :=  l.GetUUidData(uuid)
    if len(uuid_data)>0 {
        node , _ := uuid_data["node"]
        node_data := l.GetNodeData(node)
        if len(node_data) >0 {
        	 client := redis.GetRedis()
            hsetMap := make(map[string]interface{})
            hsetMap["uuid"] = uuid
            hsetMap["type"] = "msg"
            hsetMap["data"] = data
            hsetMapStringByte , _:=  json.Marshal(hsetMap)
            hsetMapString := string(hsetMapStringByte)
		     queue , _ :=  node_data["queue"].(string)
            client.Publish(queue,hsetMapString)
            l.UpdateUUidDataTime(uuid)
            return true
        }
        return false
    }
    return false
}

func (l * websocketBase) Disconnect(uuid string) bool{
	uuid_data :=  l.GetUUidData(uuid)
	if len(uuid_data)>0 {
		node , _ := uuid_data["node"]
		node_data := l.GetNodeData(node)
		if len(node_data) >0 {
			client := redis.GetRedis()
			hsetMap := make(map[string]interface{})
			hsetMap["uuid"] = uuid
			hsetMap["type"] = "close"

			hsetMapStringByte , _:=  json.Marshal(hsetMap)
			hsetMapString := string(hsetMapStringByte)
			queue , _ :=  node_data["queue"].(string)
			client.Publish(queue,hsetMapString)
 			return true
		}
		return false
	}
	return false
}


func (l * websocketBase) CloseByUUid(uuid string) bool {
 	uuid_data := l.GetUUidData(uuid)
	if len(uuid_data)>0 {
		client := redis.GetRedis()
		client.Del(fmt.Sprintf("%s%s" ,WEBSOCKET_CLIENT_HSET , uuid))
		//SwoftRedis::connection(BaseData::SESSION_REDIS)->del(self::WEBSOCKET_CLIENT_HSET.$uuid);
		//SwoftRedis::connection(BaseData::SESSION_REDIS)->hDel(self::WEBSOCKET_NODE_FD_HSET . '_' . self::getNodeId(),(string) $uuid_data['fd']);
		//SwoftRedis::connection(BaseData::SESSION_REDIS)->zRem(self::WEBSOCKET_CLIENT_ZADD,$uuid);
		client.ZRem(WEBSOCKET_CLIENT_ZADD,uuid)
		return true
 	}
	return false
}