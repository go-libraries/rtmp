package handshake

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

type Handshake struct {
	read  *bufio.Reader
	write *bufio.Writer
	Conn  *net.TCPConn
	err   error
}

func GetHandShake(tcp *net.TCPConn) *Handshake {
	return &Handshake{
		read:  bufio.NewReader(tcp),
		write: bufio.NewWriter(tcp),
		Conn:  tcp,
		err:   nil,
	}
}

//do handShake
func (hand *Handshake) DoHandshakeClient(timeout time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			hand.err = r.(error)
		}
	}()

	var err error
	//1. 发送c0
	err = hand.write.WriteByte(0x03)
	if err != nil {
		hand.setError("send c0", err)
		return
	}
	//2. 发送c1
	btC1, c1Offset := hand.buildC1Data()
	_, err = hand.write.Write(btC1)
	if err != nil {
		hand.setError("send c1", err)
		return
	}
	hand.addTimeOut(timeout)
	_ = hand.write.Flush()

	//3. 读取s0
	s0, err := hand.read.ReadByte()
	if err != nil {
		hand.setError("read s0", err)
		return
	}
	if s0 != 0x03 {
		hand.setError("check s0", fmt.Errorf("handshare Got S0: %x", s0))
		return
	}

	//4. 读取s1
	s1 := make([]byte, RtmpSigSize)
	hand.addTimeOut(timeout)
	_, err = io.ReadAtLeast(hand.read, s1, RtmpSigSize)
	if err != nil {
		hand.setError("read s1 error", err)
		return
	}
	s1Pos := hand.checkS1(s1)
	if s1Pos == 0 {
		hand.setError("check s1 error", err)
		return
	}
	//5. 读取s2 可与7 置换顺序
	s2 := make([]byte, RtmpSigSize)
	hand.addTimeOut(timeout)
	_, err = io.ReadAtLeast(hand.read, s2, RtmpSigSize)
	//6. 检测s2 对比c1 其实这个逻辑是nginx才有的，rtmp协议中规定的是随机数就行
	if err != nil || !hand.checkS2(s2, btC1, c1Offset) {
		hand.setError("check s2 error", err)
		return
	}

	//7.发送c2
	btC2 := hand.buildC2Data(s1, btC1[:4])
	_, err = hand.write.Write(btC2)
	if err != nil {
		hand.setError("send c2", err)
		return
	}
	hand.addTimeOut(timeout)
	_ = hand.write.Flush()
}

func (hand *Handshake) addTimeOut(timeout time.Duration) {
	if timeout > 0 {
		_ = hand.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}

//set error
func (hand *Handshake) setError(where string, err error) {
	if err == nil {
		hand.err = errors.New(fmt.Sprintf("%s", where))
		return
	}
	hand.err = errors.New(fmt.Sprintf("%s: %s", where, err.Error()))
}

//do handShake
func (hand *Handshake) GetError() error {
	return hand.err
}

/**
offset: 4bytes
random-data: (offset)bytes
digest-data: 32bytes
random-data: (764-4-offset-32)bytes
*/
func (hand *Handshake) buildC1Data() ([]byte, int) {
	buf := new(bytes.Buffer)

	timeBt := GetTimestampByte()
	buf.Write(timeBt)             // time (4 bytes)
	buf.Write(FlashPlayerVersion) //version (4 bytes)

	//初始化1536个随机byte
	for i := 8; i < 1536; i++ {
		buf.WriteByte(byte(rand.Int() % 256))
	}
	bts := buf.Bytes()
	//todo: 其实按照adobe 官方的，可以不用这些的步骤，只要能保证 剩下的1528个bt保持唯一就行
	// base = time(4) version(4) offset(4)
	return buildC1Digest(bts, GenuineFpKey[:30])

	/**
	Random data (1528 bytes): This field can contain any arbitrary
	 values. Since each endpoint has to distinguish between the
	 response to the handshake it has initiated and the handshake
	 initiated by its peer,this data SHOULD send something sufficiently
	 random. But there is no need for cryptographically-secure
	 randomness, or even dynamic values.
	*/
	//计算offset
	//通过随机bt 计算出 Digest的 偏移量  base is len(time)+len(version)+len(offset)
	//这里偏移量是8  4位time 4位version
	//offset := calOffset(bts, 8, 12)
	//tmpBuf := new(bytes.Buffer)
	//tmpBuf.Write(bts[:offset])
	//tmpBuf.Write(bts[offset+RtmpSha256DigestLength:])
	//
	////计算hash  内容为偏移量后的之前的数据 和 偏移量之后+32位的所有 byte
	////客户端的key 为 rtmp 协议内容
	//tempHash, err := HMACSha256(tmpBuf.Bytes(), GenuineFpKey[:30])
	//if err != nil {
	//	hand.setError("hash cal error", err)
	//	return nil, 0
	//}
	//
	////将计算结果填充到 offset后32位
	//copy(bts[offset:], tempHash)
	//return bts, offset

}

func (hand *Handshake) buildC2Data(s1 []byte, c1time []byte) []byte {
	buf := new(bytes.Buffer)

	buf.Write(c1time) //time (4 bytes)
	buf.Write(s1[:4]) //time1 (4 bytes)

	//初始化1536个随机byte
	for i := 8; i < 1536; i++ {
		buf.WriteByte(byte(rand.Int() % 256))
	}
	bts := buf.Bytes()

	s1Pos := getS1SeverPost(s1)
	digestS1, err := HMACSha256(s1[s1Pos:s1Pos+RtmpSha256DigestLength], GenuineFpKey)
	if err != nil {
		hand.setError("build c1 error:", err)
		return nil
	}
	c2Digest, err := HMACSha256(bts[:RtmpSigSize-RtmpSha256DigestLength], digestS1)
	copy(bts[RtmpSigSize-RtmpSha256DigestLength:], c2Digest)
	return bts
	//for index, b := range c2Digest {
	//	bts[RtmpSigSize-RtmpSha256DigestLength+index] = b
	//}
	//todo: 其实按照adobe 官方的，可以不用下面的步骤，只要能保证 剩下的1528个bt保持唯一就行
	/**
	Random data (1528 bytes): This field can contain any arbitrary
	 values. Since each endpoint has to distinguish between the
	 response to the handshake it has initiated and the handshake
	 initiated by its peer,this data SHOULD send something sufficiently
	 random. But there is no need for cryptographically-secure
	 randomness, or even dynamic values.
	*/
	//计算offset
	//通过随机bt 计算出 Digest的 偏移量  base is len(time)+len(version)+len(offset)
	//offset := calOffset(bts, 8, 12)
	//tmpBuf := new(bytes.Buffer)
	//tmpBuf.Write(bts[:offset])
	//tmpBuf.Write(bts[offset+RtmpSha256DigestLength:])
	//
	////计算hash  内容为偏移量后的之前的数据 和 偏移量之后+32位的所有 byte
	////客户端的key 为 rtmp 协议内容
	//tempHash, err := HMACSha256(tmpBuf.Bytes(), GenuineFpKey[:30])
	//if err != nil {
	//	hand.setError("hash cal error", err)
	//	return nil
	//}
	//
	////将计算结果填充到 offset后32位
	//copy(bts[offset:], tempHash)
	//return bts
}

func (hand *Handshake) checkS2(s2 []byte, c1 []byte, c1Offset int) bool {
	//之前的c1偏移量
	c1Hash := c1[c1Offset : c1Offset+RtmpSha256DigestLength]
	digest, err := HMACSha256(c1Hash, GenuineFmsKey)
	if err != nil {
		return false
	}

	//生成hash值 混淆是 c1的hash值 32位
	signature, err := HMACSha256(s2[:RtmpSigSize-RtmpSha256DigestLength], digest)
	if err != nil {
		return false
	}

	//检测s2内容是否与hash出来的值一致
	s2Hash := s2[RtmpSigSize-RtmpSha256DigestLength:]
	if bytes.Compare(signature, s2Hash) != 0 {
		return false
	}

	return true
}

func (hand *Handshake) checkS1(buf []byte) int {
	digestPos := calcDigestPos(buf, 8, 728)
	// Create temp buffer
	pos := checkS1(buf, digestPos)
	if pos == 0 {
		// 772 = 764+8
		digestPos = calcDigestPos(buf, 772, 728)
		return checkS1(buf, digestPos)
	}
	return pos
}

func checkS1(buf []byte, digestPos int) int {
	tmpBuf := new(bytes.Buffer)
	tmpBuf.Write(buf[:digestPos])
	tmpBuf.Write(buf[digestPos+RtmpSha256DigestLength:])
	// Generate the hash
	tempHash, err := HMACSha256(tmpBuf.Bytes(), GenuineFmsKey[:36])
	if err != nil {
		return 0
	}
	if bytes.Compare(tempHash, buf[digestPos:digestPos+RtmpSha256DigestLength]) == 0 {
		return digestPos
	}
	return 0
}
