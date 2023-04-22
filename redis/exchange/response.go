package exchange

import "goredis/pkg/utils"

var (
	pongBytes     = utils.StringToBytes("+PONG\r\n")
	oneOkResponse = new(OkResponse)
	okBytes       = utils.StringToBytes("+OK\r\n")
)

// PongResponse +PONG
type PongResponse struct{}

func (r *PongResponse) Info() []byte {
	return pongBytes
}

// OkResponse +OK
type OkResponse struct{}

func (r *OkResponse) Info() []byte {
	return okBytes
}

func MakeOkReply() *OkResponse {
	return oneOkResponse
}
