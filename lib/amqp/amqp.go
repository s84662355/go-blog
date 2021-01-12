package amqp

import (
    "github.com/streadway/amqp"
    "GoBlog/config"
    "fmt"
)

func GetAmqp() *amqp.Connection{

	  conn, err := amqp.Dial(config.AmqpConfig)
	  if err!=nil {
	  	 fmt.Println(err)
	  	 return nil
	  }
      return conn
}