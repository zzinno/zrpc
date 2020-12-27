package server

import (
	"errors"
	"github.com/vmihailenco/msgpack"
	"github.com/zzinno/zrpc/client"
	"github.com/zzinno/zrpc/logger"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

type Broker struct {
	addr             string
	port             int
	name             string
	Producers        map[string]client.Producer
	Consumer         *client.Consumer
	producersSupport map[string]client.Index
	producersClient  map[string]*rpc.Client
	Logger           logger.Logger
}

const (
	defaultAddr = "0.0.0.0"
	defaultPort = 5825
)

// @title    New
// @description  Create a service to send and receive information
// @auth      loveward         2020-12-27 12:06
// @param     addr        string        the
// @return    返回参数名        参数类型         "解释"
//addr string, port uint, producers map[string]client.Producer, consumer *client.Consumer
func (b *Broker) New(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			b.port = arg.(int)
		case string:
			b.addr = arg.(string)
		case map[string]client.Producer:
			b.Producers = arg.(map[string]client.Producer)
		case *client.Consumer:
			b.Consumer = arg.(*client.Consumer)
		case logger.Logger:
			b.Logger = arg.(logger.Logger)
		default:
			log.Fatal(errors.New("参数类型错误，请参demo中代码"))
		}
	}

	if b.addr == "" {
		b.addr = defaultAddr
	}
	if b.port == 0 {
		b.port = defaultPort
	}
	if b.Logger == nil {
		b.Logger = new(logger.ZrpcLogger)
	}

	// register

	r := rpc.NewServer()
	if b.Consumer != nil {
		err := r.Register(b.Consumer)
		checkError(err)
	}

	ipaddr, err := net.ResolveTCPAddr("tcp4", b.addr+":"+strconv.Itoa(int(b.port)))
	checkError(err)
	var listen *net.TCPListener
	listen, err = net.ListenTCP("tcp", ipaddr)
	checkError(err)
	//开启rpc
	go func(r *rpc.Server) {
		for {
			conn, Arri := listen.Accept()
			b.LogErr(Arri)
			go r.ServeConn(conn)
		}
	}(r)
	if b.Producers != nil {
		b.producersClient = make(map[string]*rpc.Client)
		b.producersSupport = make(map[string]client.Index)
		for k, v := range b.Producers {
			rpcClient, err := rpc.Dial("tcp", v.Addr+":"+strconv.Itoa(int(v.Port)))
			//rpcClient, err := rpc.DialHTTP("tcp", v.Addr+":"+strconv.Itoa(int(v.Port)))
			checkError(err)
			b.producersClient[k] = rpcClient
			ret := new([]byte)
			p := client.Params{FunName: "", Data: nil}
			err = b.producersClient[k].Call("Consumer.GetIndex", &p, ret)
			checkError(err)
			var index client.Index
			err = msgpack.Unmarshal(*ret, &index)
			checkError(err)
			b.producersSupport[k] = index
			b.Logger.Info("Function Index From ", k, ":", index.Index)
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
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func (b *Broker) LogErr(err error) {
	if err != nil {
		b.Logger.Error(err)
	}
}
