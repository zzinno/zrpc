package server

import (
	"errors"
	"github.com/vmihailenco/msgpack"
	"github.com/zzinno/zrpc/client"
	"github.com/zzinno/zrpc/logger"
	"log"
	"net"
	"net/rpc"
	"regexp"
	"strconv"
)

type Broker struct {
	addr             string
	port             int
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

/*
！！请直接在demo中查看创建的例子
@title    Broker.New
@description  创建Broker
@auth      loveward         2020-12-27 12:06
@param     addr             string                       监听地址
@param     port             int                          端口
@param     producers        map[string]client.Producer   生产者集合
@param     consumer         *client.Consumer             消费者
@param     Logger           logger.Logger                日志
*/
func (b *Broker) New(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			b.port = arg.(int)
			if b.port > 65535 || b.port < 1 {
				checkError(errors.New("端口错误：端口范围是1~65535"))
			}
		case string:
			b.addr = arg.(string)
			if !checkIp(b.addr) {
				checkError(errors.New("ip地址格式错误：给出的地址格式应该是 eg.. \"1.2.3.4\""))
			}
		case map[string]client.Producer:
			b.Producers = arg.(map[string]client.Producer)
		case *client.Consumer:
			b.Consumer = arg.(*client.Consumer)
		case logger.Logger:
			b.Logger = arg.(logger.Logger)
		default:
			log.Fatal(errors.New("参数类型错误，请参阅demo中代码"))
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

	ipaddr, err := net.ResolveTCPAddr("tcp4", b.addr+":"+strconv.Itoa(b.port))
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

			if v.Port > 65535 || v.Port < 1 {
				checkError(errors.New("端口错误：端口范围是1~65535"))
			}
			if !checkIp(v.Addr) {
				checkError(errors.New("ip地址格式错误：给出的地址格式应该是 eg.. \"1.2.3.4\""))
			}

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

func checkIp(addr string) bool {
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
	match, _ := regexp.MatchString(ipReg, addr)
	if match {
		return true
	} else {
		return false
	}
}
