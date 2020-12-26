package main

import (
	"github.com/zzinno/zrpc/client"
	"github.com/zzinno/zrpc/server"
)

func main() {
	var s server.Broker
	c := new(client.Consumer)
	c.Funs = map[string]func(*[]byte) ([]byte, error){
		"s": func(i *[]byte) ([]byte, error) {
			s := []byte("")
			return s, nil
		},
	}
	s.New("127.0.0.1", 48080, "hellos", nil, c)
	select {}
}
