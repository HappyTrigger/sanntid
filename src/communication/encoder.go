package Communication


import (
	"encoding/json"
	//"errors"
	"fmt"

)



func Encoder(message NewOrder)([]byte, error){
	return json.Marshal(message)
}


func Decoder(data []byte, message* NewOrder) {
		if err := json.Unmarshal(data, message); err!=nil{
			fmt.Println("Something went wrong when decoding", err)
		}
	return
}




