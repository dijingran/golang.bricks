package main

import (
	"math/rand"
	"os"
	"io"
	"bufio"
	"net"
	"log"
	"strconv"
	"time"
	conf "bricks/conf"
	"runtime"
)

/*
 * 题目叫码农来搬砖，题目的要求为：
 * <li>实现一个客户端和一个服务器端，客户端把服务器端的一个1G的文件搬到客户端，文件中的每行字符串为一块砖头，字符串由随机的ascii32
 * -127的字符组成，每行的长度为随机的1-200字节；
 * <li>服务端必须是单进程并且只能监听一个端口；
 * <li>客户端每个发起获取砖头的请求“线程“必须是一问一答，类似这样： 砖头 result = client.getZhuanTou();
 * 服务端收到请求后，必须是按顺序获取下一块砖头，类似这样： 砖头 Next = server.getNext();
 * <li>不允许批量处理请求，也不允许服务端处理请求的”线程“批量返回砖头，服务端处理请求的”线程“每次只能处理一个请求，并返回一块砖头；
 * <li>服务端或客户端需要对砖头进行处理
 * ，处理方式为去掉行中间的三分之一字符（从size/3字符开始去掉size/3个字符，除法向下取整)后将剩余部分以倒序的方式传输 例如 123456789
 * => 123789 => 987321 每块砖头需要标上序号，例如上面的123456789是第5行，那么最后的砖头结果应该为：5987321；
 * <li>
 * 客户端需要顺序的输出最终处理过的砖头内容到一个文件中
 * ，此文件中的砖头的顺序要和服务端的原始文件完全一致，文件不需要写透磁盘（例如java里就是可以不强制调sync）；
 * <li>不允许采用内核层面的patch；
 * <li>不建议采用通信框架，例如netty之类的这种； 不限语言、通信协议和连接数。
 * <li>比赛的运行方式： 服务端启动5s后启动客户端，客户端启动就开始计时，一直到客户端搬完所有砖头并退出计算为耗时时间；
 */
func main() {
	//	aaa := ""
	//	fmt.Println(prepare(aaa, len(aaa)))
	//	return
	log.Printf("%d", runtime.NumCPU())
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
		conn.Write([]byte{byte(gn)})// 最大127
		break;
	}

	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		if i >= len(stacks) {
			conn.Close()
			continue;
		}
		h := Handler{stacks[i], 0, len(stacks[i])}
		go h.HandleClient(conn)
		i++
	}
	//	makeFile(1024 * 1024 * 1024);
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
			break;
		}
		req = make([]byte, 200)
	}
}

// split all bricks into multi stacks of bricks.
// size : lines of each stack.
func split(size int) (a [][]string) {
	start := time.Now().Unix()
	f, err := os.Open(conf.ServerFile())
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
			line = line[0:len(line)-1]
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
	log.Printf("split into %d stacks, cost %d seconds.\n", len(stacks), (time.Now().Unix()-start))
	return stacks
}

func prepare(line string, length int) (s string) {
	if length <= 1 {
		return line
	}
	b := length / 3
//	return sort.SelectSort(line[:b] + line[b * 2:])
	return line[:b] + line[b * 2:]
}

func makeFile(max int) {
	start := time.Now().Unix();
	j := 0
	b := []byte{};
F:
	for { //
		len := rand.Intn(200)
		for i := 0; i < len ; i++ {
			j++
			if j >= max {
				break F;
			}
			b = append(b, byte(32 + rand.Intn(95)))
		}
		b = append(b, '\n')
		j++
	}
	log.Printf("Make file cost %d seconds.", (time.Now().Unix()-start))
	start = time.Now().Unix();
	f, _ := os.Create(conf.ServerFile())
	f.Write(b)
	log.Printf("Write into disk cost %d seconds.", (time.Now().Unix()-start))
}

