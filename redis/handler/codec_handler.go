package handler

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"goredis/interface/tcp"
	"goredis/pkg/log"
	"goredis/pool"
	"goredis/redis/exchange"
	"io"
	"strings"
)

func RedisCodec() codec.Codec {
	return &codecHandler{}
}

type codecHandler struct{}

func (*codecHandler) CodecName() string {
	return "codec-handler"
}

var (
	//暂不支持的命令/未知命令
	unknownOperation = []byte("-ERR unknown\r\n")
)

func (*codecHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	ch := make(chan *tcp.Request)
	//TODO 配置化协程池大小
	pools, err := pool.GetInstance(0)
	if err != nil {
		log.Errorf("run parse message with any errors, func exit: ", err)
		//TODO: 需要包装redis客户端连接关闭操作
		ctx.Channel().Close(err)
		return
	}
	err = pools.Pool.Submit(func() {
		parse(message, ch)
	})
	if err != nil {
		log.Errorf("run parse message with any errors, func exit: ", err)
		ctx.Channel().Close(err)
		return
	}
	//循环结果
	for req := range ch {
		if req.Error != nil {
			if req.Error == io.EOF ||
				req.Error == io.ErrUnexpectedEOF ||
				strings.Contains(req.Error.Error(), "use a closed network channel") {
				log.Info("handle message with errors, channel will be closed: " + ctx.Channel().RemoteAddr())
				ctx.Channel().Close(req.Error)
				return
			}
			errReply := exchange.NewStatusRequest(req.Error.Error())
			ctx.Write(errReply.RequestInfo())
			continue
		}
		if req.Data == nil {
			log.Error("empty commands")
			continue
		}
		_, ok := req.Data.(*exchange.MultiBulkRequest)
		if !ok {
			log.Error("error from multi bulk exchange")
			continue
		}
		//命令处理
		//result := h.db.Exec(message, r.Args)
		//if result != nil {
		//	ctx.Write(result.ToBytes())
		//} else {
		//	ctx.Write(unknownOperation)
		//}
	}
}

func (*codecHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	switch s := message.(type) {
	case string:
		ctx.HandleWrite(strings.NewReader(s))
	default:
		ctx.HandleWrite(message)
	}
}
