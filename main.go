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
server:= NewApiServer(":8080",store)
server.Run()
}