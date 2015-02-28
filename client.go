package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"golang.bricks/conf"
)

type Counter struct {
	count int64
}

func (c *Counter) Incr() {
	atomic.AddInt64(&c.count, 1)
}

const (
	SERVICE = "10.168.8.31:1200"
)

func main() {
	runtime.GOMAXPROCS(4)
	counter := Counter{}
	go func() {
		for {
			before := counter.count
			time.Sleep(5 * time.Second)
			log.Printf("QPS : %d, transferd %d lines.", (counter.count-before)/5, counter.count)
		}
	}()

	start := time.Now().Unix()
	ch := make(chan []string)
	gn := requestStackNum()
	log.Printf("stack num is %d .\n", gn)
	i := 0
	for ; i < gn; i++ {
		go transfer(ch, &counter)
	}
	stacks := make([][]string, gn)
	for ; i > 0; i-- {
		bricks := <-ch
		num, _ := strconv.ParseInt(bricks[0], 10, 16) // first line is stack order num, begin with 0
		if len(bricks) > 1 {
			stacks[num] = bricks[1:len(bricks)]
		} else {
			stacks[num] = []string{}
		}
	}
	log.Printf("got all bricks cost %d seconds.\n", (time.Now().Unix() - start))
	start = time.Now().Unix()
	// write into file
	f, _ := os.Create(conf.ClientFile())
	defer f.Close()
	for _, e := range stacks {
		for _, l := range e {
			f.Write([]byte(l + "\n"))
		}
	}
	log.Printf("write into disk cost %d seconds.\n", (time.Now().Unix() - start))
}

func requestStackNum() (num int) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", SERVICE)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr) // default no delay
	defer conn.Close()
	conn.Write([]byte{byte(1)})
	resp := make([]byte, 4)
	conn.Read(resp)
	return int(binary.LittleEndian.Uint16(resp))
}

func transfer(ch chan []string, counter *Counter) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", SERVICE)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr) // default no delay
	defer conn.Close()
	bricks := []string{}
	start := time.Now().Unix()
	req := []byte{byte(0)}
	emp := conf.Empty()
	fin := conf.Finish()
	for {
		conn.Write(req)
		resp := make([]byte, 200)
		len, _ := conn.Read(resp)
		s := string(resp[0:len])
		if s == emp {
			s = ""
		}
		if s == fin {
			break
		}
		bricks = append(bricks, s)
		counter.Incr()
	}
	ch <- bricks
	fmt.Printf("finished, cost %d seconds.\n", time.Now().Unix()-start)
}
