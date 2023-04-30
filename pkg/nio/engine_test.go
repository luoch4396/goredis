package nio

import (
	log2 "goredis/pkg/log"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var addr = "127.0.0.1:8848"
var testfile = "test_nio_server.file"

func init() {
	fs := &log2.FileSettings{
		Path:     "logs",
		FileName: "goredis",
	}
	//初始化日志模块
	log2.NewLoggerBuilder().
		BuildStdOut(os.Stdout).
		BuildLevel("DEBUG").
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
		c.SetKeepAlive(true, time.Second*60)
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
	g.OnReadBufferFree(func(c *Conn, b []byte) {

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
