package main

import (
	"github.com/ZZINNO/zrpc/client"
	"github.com/ZZINNO/zrpc/server"
)

func main() {
	var s server.Broker
	c := new(client.Consumer)
	c.Funs = map[string]func(*[]byte) ([]byte, error){
		"say": func(i *[]byte) ([]byte, error) {
			s := []byte("world")
			return s, nil
		},
	}
	s.New("127.0.0.1", 48080, "hellos", nil, c)
	select {}
}
