//go:build darwin || netbsd || freebsd || openbsd || dragonfly

package nio

import (
	"goredis/pkg/log"
	"goredis/pkg/pool/bytepool"
	"io"
	"os"
)

var bufSize = 1024 * 16

func (c *Conn) Sendfile(f *os.File, remain int64) (written int64, err error) {
	if f == nil {
		return 0, nil
	}

	if remain <= 0 {
		stat, err := f.Stat()
		if err != nil {
			return 0, err
		}
		remain = stat.Size()
	}

	for remain > 0 {
		if bufSize > int(remain) {
			bufSize = int(remain)
		}
		buf := bytepool.Malloc(bufSize)
		nr, er := f.Read(buf)
		if nr > 0 {
			nw, ew := c.Write(buf[0:nr])
			if nw < 0 {
				nw = 0
			}
			remain -= int64(nw)
			written += int64(nw)
			if ew != nil {
				err = ew
				log.Errorf("Sendfile error: ", err)
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				log.Errorf("Sendfile error: ", err)
				break
			}
		}
		if er != nil && er != io.EOF {
			err = er
			break
		}
	}
	return written, err
}
