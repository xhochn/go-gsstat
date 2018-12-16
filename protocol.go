package gsstat

import (
	"net"
	"time"
)

// RequestUDP ..
type RequestUDP struct {
	Conn   *net.UDPConn
	Buf    []byte
	Offset int
}

// NewUDPConnection ..
func (r *RequestUDP) NewUDPConnection(addr string, timeout time.Duration) (err error) {
	fullAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}
	r.Conn, err = net.DialUDP("udp", nil, fullAddr)
	if err != nil {
		return
	}
	if err = r.Conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return
	}
	return
}

// Close ..
func (r *RequestUDP) Close() {
	r.Conn.Close()
}

// Send ..
func (r *RequestUDP) Send(raw []byte) (err error) {
	_, err = r.Conn.Write(raw)
	return
}

// ReadFrom ..
func (r *RequestUDP) ReadFrom() (err error) {
	buf := make([]byte, 1400)
	l, _, err := r.Conn.ReadFrom(buf)
	if err != nil {
		return
	}
	r.Buf = buf[:l]
	return
}

// SetOffset ..
func (r *RequestUDP) SetOffset(n int) {
	r.Offset = n
}

// Read ..
func (r *RequestUDP) Read(count int) []byte {
	b := r.Buf[r.Offset : r.Offset+count]
	r.Offset += count
	return b
}

func (r *RequestUDP) String() string {
	// look for \0
	soffset := r.Offset
	for {
		if r.Read(1)[0] == 0 {
			break
		}
	}
	return string(r.Buf[soffset : r.Offset-1])
}
