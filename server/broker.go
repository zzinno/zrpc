package server

import (
	"errors"
	"fmt"
	"github.com/ZZINNO/zrpc/client"
	"github.com/vmihailenco/msgpack"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

type Broker struct {
	addr             string
	port             uint
	name             string
	Producers        map[string]client.Producer
	Consumer         *client.Consumer
	producersSupport map[string]client.Index
	producersClient  map[string]*rpc.Client
}

func (b *Broker) New(addr string, port uint, name string, producers map[string]client.Producer, consumer *client.Consumer) {
	// init
	b.addr = addr
	b.port = port
	b.name = name
	b.Producers = producers
	b.Consumer = consumer
	// register
	_ = rpc.Register(consumer)
	rpc.HandleHTTP()

	//开启rpc
	go func() {
		err := http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	if producers != nil {
		b.producersClient = make(map[string]*rpc.Client)
		b.producersSupport = make(map[string]client.Index)
		for k, v := range b.Producers {
			rpcClient, err := rpc.DialHTTP("tcp", v.Addr+":"+strconv.Itoa(int(v.Port)))
			if err != nil {
				log.Fatal(err)
			}
			b.producersClient[k] = rpcClient
			ret := new([]byte)
			p := client.Params{FunName: "", Data: nil}
			err = b.producersClient[k].Call("Consumer.GetIndex", &p, ret)
			if err != nil {
				log.Fatal(err)
			}
			var index client.Index
			err = msgpack.Unmarshal(*ret, &index)
			if err != nil {
				log.Fatal(err)
			}
			b.producersSupport[k] = index
			fmt.Println("Function Index From ", k, ":", index.Index)
		}
	}
}

func (b *Broker) Call(producerName string, funName string, data *[]byte) (*[]byte, error) {
	if b.checkFunSafe(producerName, funName) {
		ret := new([]byte)
		p := client.Params{FunName: funName, Data: data}
		err := b.producersClient[producerName].Call("Consumer.Deal", &p, ret)
		return ret, err
	} else {
		return nil, errors.New("Function Name Not Find ")
	}

}
func (b *Broker) checkFunSafe(producerName string, funName string) bool {
	for _, v := range b.producersSupport[producerName].Index {
		if funName == v {
			return true
		}
	}
	return false
}
