package server

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
)

func FastCollect(ctx *fasthttp.RequestCtx) {
	//go func(bts []byte) {
	//fmt.Println(string(ctx.Path()))
	uris := bytes.Split(ctx.Path(), []byte("/"))
	aId, err := strconv.ParseInt(string(uris[3]), 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	err = c.Collect(string(uris[2]), aId)
	if err != nil {
		fmt.Println(err)
	}
	//}(ctx.Path())

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Content-Length", "4")
	ctx.Write([]byte("true"))
}
