package utilities


import (
	"encoding/json"
	"log"



)

func Encoder(message Message)([]byte){
	result, err :=json.Marshal(message)
	if err!=nil{
		log.Fatal("Encoder failed",err)
	}
	return result
}


func Decoder(data []byte) Message {
		var message Message
		if err := json.Unmarshal(data, message); err!=nil{
			log.Fatal("Decoder failed", err)
		}
	return message
}




