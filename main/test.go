package main

import (
	"fmt"
	"time"
	/// "sync"
	"GoBlog/lib/nsqclient"
	_ "errors"
	"github.com/nsqio/go-nsq"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(10)
}

func main() {

	ttt := nsqclient.ConnectNsqProducer("127.0.0.1:4150", "safds")

	ttt.Publish("dcsmcls", []byte("dfmvfdlkvmdfkl"))

	nsqclient.NsqConsumer("dsdcds", "dcsmcls", "11111111", []string{"127.0.0.1:4161"}, func(message *nsq.Message) error {
		//	time.Sleep(1 * time.Second)
		fmt.Println(message.Attempts)
		fmt.Println(string(message.Body))

		///	message.Finish()
		message.Requeue(1)
		fmt.Println("111111111")
		return nil
		//return errors.New("dsmlvmdkvld")
	}, 1, 3)

	time.Sleep(400 * time.Second)

}
