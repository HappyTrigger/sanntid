package utilities


import (
	"encoding/json"
	//"errors"
	"fmt"

)

func Encoder(message Message)([]byte, error){
	return json.Marshal(message)
}


func Decoder(data []byte, message* Message) {
		if err := json.Unmarshal(data, message); err!=nil{
			fmt.Println("Something went wrong when decoding", err)
		}
	return
}




