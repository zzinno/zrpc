package main

import (
	"fmt"
	"github.com/zzinno/zrpc/client"
	"github.com/zzinno/zrpc/server"
	"sync"
	"time"
)

var b server.Broker
var a = []byte("")
var errNum int

func main() {
	test(38480)
	errNum = 0
	test(38481)
	errNum = 0
	test(38482)
	errNum = 0
	test(38483)
}
func sendMessage(du *chan int64) error {
	start1 := time.Now().UnixNano() / int64(time.Millisecond)

	_, err := b.Call("hello", "s", &a)
	//fmt.Println(re.Data)
	end1 := time.Now().UnixNano() / int64(time.Millisecond)
	if err != nil {
		errNum += 1
		return err
	}
	cost := end1 - start1
	*du <- cost
	return nil
}

func test(port int) {
	testNum := 500000
	du := make(chan int64, testNum)
	var wg sync.WaitGroup
	start := time.Now().UnixNano() / int64(time.Millisecond)

	// 处理连接

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
	b.New("127.0.0.1", port, p, c)

	var max int64
	var test int64
	for {
		wg.Add(1)
		test += 1
		go func(wg *sync.WaitGroup, du *chan int64) {
			defer wg.Done()
			err := sendMessage(du)
			if err != nil {
				return
			}
		}(&wg, &du)
		if test == int64(testNum) {
			break
		}
	}
	wg.Wait()
	end := time.Now().UnixNano() / int64(time.Millisecond)
	lent := len(du)
	count := len(du)
	var allcost int64
	for v := range du {
		lent -= 1
		if max < v {
			max = v
		}
		allcost += v
		if lent <= 0 {
			break
		}
	}
	fmt.Println("request count:", testNum)
	fmt.Println("success rate:", count/testNum*100, "%")
	fmt.Println("max cost:", max, "ms")
	fmt.Println("avg cost:", allcost/test, "ms")
	fmt.Println("all cost:", end-start, "ms")
	fmt.Println("ERRNUM:", errNum)
	fmt.Println("rps:", int64(testNum)*1000/(end-start), "request/s")
}
