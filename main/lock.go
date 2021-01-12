package main
import (
    "os"
    "os/signal"
    "syscall"

  ///  "GoBlog/lib/helper"
  ///  "fmt"
     "GoBlog/lib/log"
)


func init(){
     log.Logger.Init("log/lock",true)
}

func main() {
    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)


    
 
    log.Logger.Info("384954398589")

   //  process.test()

    <-done
   // logger.Sugar.Info("process exit right now")


}