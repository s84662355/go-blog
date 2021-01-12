package event

import (
  "github.com/goinggo/mapstructure"
)

func GetParam(data map[string]interface{}, body interface{}){
        mapstructure.DecodePath(data,body)
}