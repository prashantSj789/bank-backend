package main

import (
	"fmt"
	"log"
)



func main(){
	store,err:= NewPostgressStore()
	if err!=nil{
		log.Fatal("error")
	}
	if er,err,err1,err2:= store.init();err!=nil {
	log.Fatal("Error",err)
	log.Fatal("Error",er)
	log.Fatal("Error",err1)
	log.Fatal("Error",err2)
	}
server:= NewApiServer(":8080",store)
server.Run()
 fmt.Println("hi there")
}