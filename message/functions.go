package message

import (
	"bufio"

	"io"
	"net"
	"time"
)

// Get timestamp
func GetTimestamp() uint32 {
	//return uint32(0)
	return uint32(time.Now().UnixNano()/int64(1000000)) % MaxTimestamp
}

// Read byte from network
func ReadByteFromNetwork(r Reader) (b byte, err error) {
	retry := 1
	for {
		b, err = r.ReadByte()
		if err == nil {
			return
		}
		netErr, ok := err.(net.Error)
		if !ok {
			return
		}
		if !netErr.Temporary() {
			return
		}

		if retry < 16 {
			retry = retry * 2
		}
		time.Sleep(time.Duration(retry*100) * time.Millisecond)
	}
	return
}

// Read bytes from network
func ReadAtLeastFromNetwork(r Reader, buf []byte, min int) (n int, err error) {
	retry := 1
	for {
		n, err = io.ReadAtLeast(r, buf, min)
		if err == nil {
			return
		}
		netErr, ok := err.(net.Error)
		if !ok {
			return
		}
		if !netErr.Temporary() {
			return
		}

		if retry < 16 {
			retry = retry * 2
		}
		time.Sleep(time.Duration(retry*100) * time.Millisecond)
	}
	return
}

// Copy bytes from network
func CopyNFromNetwork(dst Writer, src Reader, n int64) (written int64, err error) {
	// return io.CopyN(dst, src, n)

	buf := make([]byte, 4096)
	for written < n {
		l := len(buf)
		if d := n - written; d < int64(l) {
			l = int(d)
		}
		nr, er := ReadAtLeastFromNetwork(src, buf[0:l], l)
		if er != nil {
			err = er
			break
		}
		if nr == l {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		} else {
			err = io.ErrShortBuffer
		}
	}
	return
}

func WriteToNetwork(w Writer, data []byte) (written int, err error) {
	length := len(data)
	var n int
	retry := 1
	for written < length {
		n, err = w.Write(data[written:])
		if err == nil {
			written += int(n)
			continue
		}
		netErr, ok := err.(net.Error)
		if !ok {
			return
		}
		if !netErr.Temporary() {
			return
		}

		if retry < 16 {
			retry = retry * 2
		}
		time.Sleep(time.Duration(retry*500) * time.Millisecond)
	}
	return

}

// Copy bytes to network
func CopyNToNetwork(dst Writer, src Reader, n int64) (written int64, err error) {
	// return io.CopyN(dst, src, n)

	buf := make([]byte, 4096)
	for written < n {
		l := len(buf)
		if d := n - written; d < int64(l) {
			l = int(d)
		}
		nr, er := io.ReadAtLeast(src, buf[0:l], l)
		if nr > 0 {
			nw, ew := WriteToNetwork(dst, buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return
}

func FlushToNetwork(w *bufio.Writer) (err error) {
	retry := 1
	for {
		err = w.Flush()
		if err == nil {
			return
		}
		netErr, ok := err.(net.Error)
		if !ok {
			return
		}
		if !netErr.Temporary() {
			return
		}

		if retry < 16 {
			retry = retry * 2
		}
		time.Sleep(time.Duration(retry*500) * time.Millisecond)
	}
	return
}
