package main

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"golang.bricks/conf"
)

func main() {
	//makeFile(1024 * 1024 * 1024)
	//return
	log.Printf("runtime.NumCPU : %d", runtime.NumCPU())
	runtime.GOMAXPROCS(2)
	stacks := split(512 * 1024)
	gn := len(stacks)
	service := ":1200"
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", service)
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	// tell client the stack num
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		req := make([]byte, 4)
		conn.Read(req)
		conn.Write([]byte{byte(gn)}) // 最大127
		break
	}

	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		if i >= len(stacks) {
			conn.Close()
			continue
		}
		h := Handler{stacks[i], 0, len(stacks[i])}
		go h.HandleClient(conn)
		i++
	}
}

type Handler struct {
	data []string
	i    int
	max  int
}

func (h *Handler) Next() (s string) {
	if h.i >= h.max {
		return conf.Finish()
	}
	line := h.data[h.i]
	h.i++
	return line
}

func (h *Handler) HandleClient(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 200)
	for {
		conn.Read(req)
		line := h.Next()
		conn.Write([]byte(line)) // don't care about return value
		//		conn.Write([]byte(prepare(line,len(line)))) // don't care about return value
		if line == conf.Finish() {
			log.Println("Finished.")
			break
		}
		req = make([]byte, 200)
	}
}

// split all bricks into multi stacks of bricks.
// size : lines of each stack.
func split(size int) (a [][]string) {
	start := time.Now().Unix()
	f, err := os.Open(conf.SERVER_FILE)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	g := []string{}
	stacks := [][]string{}
	for i := 0; ; i++ {
		if i == 0 {
			g = append(g, strconv.Itoa(len(stacks))) // each stack begin with an order num
		}
		line, err := reader.ReadString(byte('\n'))
		if err != nil || io.EOF == err {
			stacks = append(stacks, g)
			break
		}
		if len(line) >= 1 { // delete '\n'
			line = line[0 : len(line)-1]
		}
		length := len(line)
		if length == 0 {
			line = conf.Empty()
		}
		//		g = append(g, line)
		g = append(g, prepare(line, length))
		if i+1 >= size {
			i = -1
			stacks = append(stacks, g)
			g = []string{}
		}
	}
	log.Printf("split into %d stacks, cost %d seconds.\n", len(stacks), (time.Now().Unix() - start))
	return stacks
}

func prepare(line string, length int) (s string) {
	if length <= 1 {
		return line
	}
	b := length / 3
	//	return sort.SelectSort(line[:b] + line[b * 2:])
	return line[:b] + line[b*2:]
}

func makeFile(max int) {
	start := time.Now().Unix()
	j := 0
	b := []byte{}
F:
	for { //
		len := rand.Intn(200)
		for i := 0; i < len; i++ {
			j++
			if j >= max {
				break F
			}
			b = append(b, byte(32+rand.Intn(95)))
		}
		b = append(b, '\n')
		j++
	}
	log.Printf("Make file cost %d seconds.", (time.Now().Unix() - start))
	start = time.Now().Unix()
	f, _ := os.Create(conf.SERVER_FILE)
	f.Write(b)
	log.Printf("Write into disk cost %d seconds.", (time.Now().Unix() - start))
}
