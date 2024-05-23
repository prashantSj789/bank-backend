package main

import (

	"log"
)



func main(){
	store,err:= NewPostgressStore()
	if err!=nil{
		log.Fatal("error")
	}
	if er,err:= store.init();err!=nil {
	log.Fatal("Error",err)
	log.Fatal("Error",er)
	}
server:= NewApiServer(":3000",store)
server.Run()
}