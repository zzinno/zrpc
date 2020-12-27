package main

import (
	"fmt"
	"github.com/zzinno/zrpc/client"
	"github.com/zzinno/zrpc/server"
)

func main() {
	var b server.Broker
	p := map[string]client.Producer{
		"hello": {
			Addr: "127.0.0.1",
			Port: 48080,
		},
	}
	c := new(client.Consumer)
	c.Funs = map[string]func(*[]byte) ([]byte, error){
		"s": func(i *[]byte) ([]byte, error) {
			s := []byte("world")
			return s, nil
		},
	}
	b.New("127.0.0.1", 38080, p, c)
	a := []byte("")
	ret, err := b.Call("hello", "s", &a)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(*ret), err)
	}

}
