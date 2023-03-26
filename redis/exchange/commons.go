package exchange

var (
	nullBulkBytes       = []byte("$-1\r\n")
	CRLF                = "\r\n"
	emptyMultiBulkBytes = []byte("*0\r\n")
)

// StatusInfo 状态指令
type StatusInfo struct {
	Status string
}

func NewStatusInfo(status string) *StatusInfo {
	return &StatusInfo{
		Status: status,
	}
}

func (r *StatusInfo) Info() []byte {
	return []byte("+" + r.Status + CRLF)
}
