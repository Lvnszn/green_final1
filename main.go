package main

import (
	"github.com/valyala/fasthttp"
	"green_final1/pkg/server"
	"log"
)

func main() {
	err := fasthttp.ListenAndServe(":8080", server.FastCollect)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//http.HandleFunc("/collect_energy/", server.Collect)
	//addr := ":8080"
	////if Env == "prod" {
	////	addr = ":80"
	////}
	//err := http.ListenAndServe(addr, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
