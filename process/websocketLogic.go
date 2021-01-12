package process

import (
     
      "GoBlog/lib/amqp"
    
      "GoBlog/logic/websocket/event"
      
)


type websocketLogic struct {
   
}

func (l * websocketLogic) Run(){
    
    conn := amqp.GetAmqp() 
    defer conn.Close()
 
     ch ,_  := conn.Channel()
     defer ch.Close()
           q ,_  := ch.QueueDeclare("task_queue",false,true,false,false,nil)
            ch.Qos(
                1 ,     // prefetch count
                0,     // prefetch size
                false, // global
        )
             msgs ,_ := ch.Consume(
                q.Name, // queue
                "",     // consumer
            true, false, false, false,  nil,  )  // args )

            
        for d := range msgs {
            event.Run(string(d.Body))
        }
          
}
