//go:build linux
// +build linux

package nio

func (c *Conn) Sendfile(f *os.File, remain int64) (written int64, err error) {

}
