package server

import (
	"green_final1/pkg/service"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	//c = service.NewService()
	//c = service.NewMemoryService()
}

var c service.Collector

// Collect /collect_energy/user001/energy001
func Collect(w http.ResponseWriter, request *http.Request) {
	uris := strings.Split(request.RequestURI, "/")
	aId, err := strconv.ParseInt(uris[3], 10, 64)
	if err != nil {
		io.WriteString(w, "true")
		return
	}
	err = c.Collect(uris[2], aId)
	if err != nil {
		io.WriteString(w, "true")
		return
	}
	w.WriteHeader(200)
	io.WriteString(w, "true")
}
