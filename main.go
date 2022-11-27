package main

import (
	"github.com/valyala/fasthttp"
	_ "go.uber.org/automaxprocs"
	"green_final1/pkg/server"
	"log"
)

func main() {
	//s := server.New("tcp://0.0.0.0:8080", true, func(hc *server.HttpCodec, body []byte) {
	//	hc.Suc()
	//})
	//fmt.Printf("server stopped: %v\n", s.Run())

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
