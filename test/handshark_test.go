package test

import (
	"fmt"
	"net"
	"rtmp/handshake"
	"testing"
	"time"
)

func TestHandshake(t *testing.T)  {
	//var s = "rtmp://zbds.hnradio.com/live_zx/938_120K?auth_key=1597737496-0-0-4abf1801d28f9d5f1441fa8159f5af2b"
	//s1 := &s
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "zbds.hnradio.com", 1935))
	//119.36.225.198
	//tcpConn, e := net.DialTCP("tcp", nil, &net.TCPAddr{
	//	//	IP:   net.IPv4(119, 36, 225, 198),
	//	//	Port: 1935,
	//	//})
	tcpConn := c.(*net.TCPConn)
	fmt.Println(tcpConn, err)
	hand := handshake.GetHandShake(tcpConn)
	hand.DoHandshakeClient(3*time.Second)
	fmt.Println(hand.GetError())
}
