package exchange

type PongResponse struct{}

var pongBytes = []byte("+PONG\r\n")

func (r *PongResponse) ResponseInfo() []byte {
	return pongBytes
}
