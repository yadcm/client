package main

import (
	"context"
	"log"
	"net/netip"
	"yadcmc/internal/app"
	protocol "yadcmc/internal/pb/yadcmd.daemon"
)

func main() {
	// https://github.com/gizak/termui
	cl, err := app.NewClient(netip.MustParseAddrPort("0.0.0.0:49069"))
	if err != nil {
		panic(err)
	}
	result, err := cl.Greeting(context.Background(), &protocol.Hello{Version: 1000})
	if err != nil {
		panic(err)
	}
	log.Println(result.Version)
}
