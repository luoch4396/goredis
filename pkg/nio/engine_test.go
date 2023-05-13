package nio

import (
	log2 "goredis/pkg/log"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var addr = "127.0.0.1:8888"
var testfile = "test.txt"
var e *Engine

func init() {
	//日志
	fs := &log2.FileSettings{
		Path:     "logs",
		FileName: "goredis",
	}
	//初始化日志模块
	log2.NewLoggerBuilder().
		BuildStdOut(os.Stdout).
		BuildLevel("INFO").
		BuildFile(fs).
		Build()

	if err := os.WriteFile(testfile, make([]byte, 1024*100), 0600); err != nil {
		log.Panicf("write file failed: %v", err)
	}

	addrs := []string{addr}
	g := NewEngine(Config{
		Network: "tcp",
		Addrs:   addrs,
	})

	g.OnOpen(func(c *Conn) {
		c.SetReadDeadline(time.Now().Add(time.Second * 10))
	})
	g.OnData(func(c *Conn, data []byte) {
		if len(data) == 8 && string(data) == "sendfile" {
			fd, err := os.Open(testfile)
			if err != nil {
				log.Panicf("open file failed: %v", err)
			}

			if _, err = c.Sendfile(fd, 0); err != nil {
				panic(err)
			}

			if err := fd.Close(); err != nil {
				log.Panicf("close file failed: %v", err)
			}
		} else {
			c.Write(append([]byte{}, data...))
		}
	})
	g.OnClose(func(c *Conn, err error) {})

	err := g.Start()
	if err != nil {
		log.Panicf("Start failed: %v\n", err)
	}

	e = g
}

func TestEcho(t *testing.T) {
	var done = make(chan int)
	var clientNum = 2
	var msgSize = 1024
	var total int64 = 0

	g := NewEngine(Config{})
	err := g.Start()
	if err != nil {
		log.Panicf("Start failed: %v\n", err)
	}
	defer g.Stop()

	g.OnOpen(func(c *Conn) {
		c.SetSession(1)
		if c.Session() != 1 {
			log.Panicf("invalid session: %v", c.Session())
		}
		c.SetLinger(1, 0)
		c.SetNoDelay(true)
		//c.SetKeepAlive(true, time.Second*60)
		c.SetDeadline(time.Now().Add(time.Second))
		c.SetReadBuffer(1024 * 4)
		c.SetWriteBuffer(1024 * 4)
		log.Printf("connected, local addr: %v, remote addr: %v", c.LocalAddr(), c.RemoteAddr())
	})
	g.BeforeWrite(func(c *Conn) {
		c.SetWriteDeadline(time.Now().Add(time.Second * 5))
	})
	g.OnData(func(c *Conn, data []byte) {
		recved := atomic.AddInt64(&total, int64(len(data)))
		if len(data) > 0 && recved >= int64(clientNum*msgSize) {
			close(done)
		}
	})

	g.OnReadBufferAlloc(func(c *Conn) []byte {
		return make([]byte, 1024)
	})

	one := func(n int) {
		c, err := Dial("tcp", addr)
		if err != nil {
			log.Panicf("Dial failed: %v", err)
		}
		g.AddConn(c)
		if n%2 == 0 {
			c.Writev([][]byte{make([]byte, msgSize)})
		} else {
			c.Write(make([]byte, msgSize))
		}
	}

	for i := 0; i < clientNum; i++ {
		if runtime.GOOS != "windows" {
			one(i)
		} else {
			go one(i)
		}
	}

	<-done
}

func TestSendfile(t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024*100)

	for i := 0; i < 3; i++ {
		if _, err := conn.Write([]byte("sendfile")); err != nil {
			log.Panicf("write 'sendfile' failed: %v", err)
		}

		if _, err := io.ReadFull(conn, buf); err != nil {
			log.Panicf("read file failed: %v", err)
		}
	}
}

func TestUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	unixAddr := "./test.unix"
	defer os.Remove(unixAddr)
	g := NewEngine(Config{
		Network: "unix",
		Addrs:   []string{unixAddr},
	})
	var connSvr *Conn
	var connCli *Conn
	g.OnOpen(func(c *Conn) {
		if connSvr == nil {
			connSvr = c
		}
		c.Type()
		c.IsTCP()
		c.IsUnix()
		log.Printf("unix onOpen: %v, %v", c.LocalAddr().String(), c.RemoteAddr().String())
	})
	g.OnData(func(c *Conn, data []byte) {
		log.Println("unix onData:", c.LocalAddr().String(), c.RemoteAddr().String(), string(data))
		if c == connSvr {
			_, err := c.Write([]byte("world"))
			if err != nil {
				t.Fatal(err)
			}
		}
		if c == connCli && string(data) == "world" {
			c.Close()
		}
	})
	chClose := make(chan *Conn, 2)
	g.OnClose(func(c *Conn, err error) {
		log.Println("unix onClose:", c.LocalAddr().String(), c.RemoteAddr().String(), err)
		chClose <- c
	})

	err := g.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer g.Stop()

	c, err := net.Dial("unix", unixAddr)
	if err != nil {
		t.Fatalf("unix Dial: %v, %v, %v", c.LocalAddr(), c.RemoteAddr(), err)
	}
	defer c.Close()
	time.Sleep(time.Second / 10)
	buf := []byte("hello")
	connCli, err = g.AddConn(c)
	if err != nil {
		t.Fatalf("unix AddConn: %v, %v, %v", c.LocalAddr(), c.RemoteAddr(), err)
	}
	_, err = connCli.Write(buf)
	if err != nil {
		t.Fatalf("unix Write: %v, %v, %v", c.LocalAddr(), c.RemoteAddr(), err)
	}
	<-chClose
	<-chClose
}
