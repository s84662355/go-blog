package config

import (
   // "github.com/tidwall/gjson"
   // "io/ioutil"
  //  "log"
   // "strconv"
    
)

func init() {
    /*
    content, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal(err)
    }

    configStr := string(content)

    var value string
    var values []string

    //环境
    value = gjson.Get(configStr, "env").String()
    if value != "" {
        Env = value
    }
    }

    //mysql 配置

    //mysql 配置
    value = gjson.Get(configStr, "mysql").String()
    if value != "" {
        MySQL = value
    }

    value = gjson.Get(configStr, "mysql_host").String()
    if value != "" {
        Host = value
    }
    value = gjson.Get(configStr, "mysql_port").String()
    if value != "" {
        Port = value
    }
    value = gjson.Get(configStr, "mysql_user").String()
    if value != "" {
        User = value
    }
    value = gjson.Get(configStr, "mysql_pwd").String()
    if value != "" {
        Pwd = value
    }
    value = gjson.Get(configStr, "mysql_name").String()
    if value != "" {
        Name = value
    }

    value = gjson.Get(configStr, "redis_address").String()
    if value != "" {
        RedisIP = value
    }
    value = gjson.Get(configStr, "redis_password").String()
    if value != "" {
        RedisPassword = value
    } else if Env == "local" {
        RedisPassword = value
    }

    //memcache 配置，数组
    values = make([]string, 0)
    for _, v := range gjson.Get(configStr, "memcache").Array() {
        values = append(values, v.String())
    }
    if len(values) > 0 {
        Memcache = values
    }

    //阿里云表格存储
    value = gjson.Get(configStr, "tablestore_host").String()
    if value != "" {
        TableStoreEndPoint = value
    }
    value = gjson.Get(configStr, "tablestore_db").String()
    if value != "" {
        TableStoreDB = value
    }

    //七牛
    value = gjson.Get(configStr, "qiniu_host").String()
    if value != "" {
        QiniuHost = value
    }
    value = gjson.Get(configStr, "qiniu_bucket").String()
    if value != "" {
        QiniuBucket = value
    }

    //消息队列
    values = make([]string, 0)
    for _, v := range gjson.Get(configStr, "nsq").Array() {
        values = append(values, v.String())
    }
    if len(values) > 0 {
        NSQConsumers = values
    }

    //nsq tcp hosts
    NSQServerHosts[helper.GetLocalIP()+"."+ConnectTCPListenPort] = struct{}{}
    for _, v := range gjson.Get(configStr, "nsq_server_hosts").Array() {
        NSQServerHosts[v.String()+"."+ConnectTCPListenPort] = struct{}{}
    }

    //oa
    value = gjson.Get(configStr, "oa_url").String()
    if value != "" {
        OAUrl = value
    }

    //ybs
    value = gjson.Get(configStr, "ybs_url").String()
    if value != "" {
        YbsUrl = value
    }
    value = gjson.Get(configStr, "ybs_head").String() //ybs 总部职能的部门id
    if value != "" {
        CenterEntityId, _ = strconv.Atoi(value)
    }

    value = gjson.Get(configStr, "ybs_guangzhouPianqu").String() //ybs 广州片区部门id
    if value != "" {
        GuangzhouPianqu, _ = strconv.Atoi(value)
    }

    value = gjson.Get(configStr, "jz_center_entity_id").String() //ybs 软件总部id
    if value != "" {
        JzCenterEntityId, _ = strconv.Atoi(value)
    }

    value = gjson.Get(configStr, "yz_center_entity_id").String() //ybs 伊智总部id
    if value != "" {
        YzCenterEntityId, _ = strconv.Atoi(value)
    }

    //umeng
    value = gjson.Get(configStr, "umeng_production").String() //ybs 总部职能的部门id
    if value != "" {
        UmengProductionMode = gjson.Get(configStr, "umeng_production").Bool()
    }

    value = gjson.Get(configStr, "app_domain").String() //项目域名
    if value != "" {
        AppDomain = value
    }
    */
}
