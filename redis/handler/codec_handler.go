package handler

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pool/gopool"
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
	var handleReadFunc = func() {
		handleRead(ctx, message)
	}
	gopool.Go(handleReadFunc)
}

func handleRead(ctx netty.InboundContext, message netty.Message) {
	ch := make(chan *tcp.Request)
	var parseStreamingFunc = func() {
		parseStreaming(message, ch)
	}
	gopool.Go(parseStreamingFunc)
	//循环结果
	for req := range ch {
		if req.Error != nil {
			if req.Error == io.EOF || req.Error == io.ErrUnexpectedEOF ||
				strings.Contains(req.Error.Error(), "use a closed network channel") {
				log.Errorf("handle message with errors, channel will be closed: " + ctx.Channel().RemoteAddr())
				ctx.Channel().Close(req.Error)
				return
			}
			errRep := errors.NewStandardError(req.Error.Error())
			ctx.Write(errRep.Info())
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
	//switch s := message.(type) {
	//case string:
	//	ctx.HandleWrite(strings.NewReader(s))
	//default:
	//	ctx.HandleWrite(message)
	//}
}
