package handshake

const C0S0 = 0x03

const RtmpSigSize = 1536 // time (4 bytes) version (4 bytes) key (764 bytes)   digest (764 bytes)
const RtmpLargeHeaderSize = 12
const RtmpSha256DigestLength = 32
const RtmpDefaultChunkSize = 128

const (
	MaxTimeStamp  = uint64(2000000000)
	AutoTimeStamp = uint32(0XFFFFFFFF)

	DefaultHighPriorityBufferSize   = 2048
	DefaultMiddlePriorityBufferSize = 128
	DefaultLowPriorityBufferSize    = 64
	DefaultChunkSize                = uint32(128)
	DefaultWindowSize               = 2500000
	DefaultCapabilities             = float64(15)
	DefaultAudioCode                = float64(4071)
	DefaultVideoCode                = float64(252)
	FmsCapabilities                 = uint32(255)
	FmsMode                         = uint32(2)

	SetPeerBandwidthHard    = byte(0)
	SetPeerBandwidthSoft    = byte(1)
	SetPeerBandwidthDynamic = byte(2)
)

var (
	FlashPlayerVersion = []byte{0x09, 0x00, 0x7C, 0x02}
	GenuineFmsKey      = []byte{
		0x47, 0x65, 0x6e, 0x75, 0x69, 0x6e, 0x65, 0x20,
		0x41, 0x64, 0x6f, 0x62, 0x65, 0x20, 0x46, 0x6c,
		0x61, 0x73, 0x68, 0x20, 0x4d, 0x65, 0x64, 0x69,
		0x61, 0x20, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
		0x20, 0x30, 0x30, 0x31, // Genuine Adobe Flash Media Server 001
		0xf0, 0xee, 0xc2, 0x4a, 0x80, 0x68, 0xbe, 0xe8,
		0x2e, 0x00, 0xd0, 0xd1, 0x02, 0x9e, 0x7e, 0x57,
		0x6e, 0xec, 0x5d, 0x2d, 0x29, 0x80, 0x6f, 0xab,
		0x93, 0xb8, 0xe6, 0x36, 0xcf, 0xeb, 0x31, 0xae,
	}
	GenuineFpKey = []byte{
		0x47, 0x65, 0x6E, 0x75, 0x69, 0x6E, 0x65, 0x20,
		0x41, 0x64, 0x6F, 0x62, 0x65, 0x20, 0x46, 0x6C,
		0x61, 0x73, 0x68, 0x20, 0x50, 0x6C, 0x61, 0x79,
		0x65, 0x72, 0x20, 0x30, 0x30, 0x31, /* Genuine Adobe Flash Player 001 */
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8,
		0x2E, 0x00, 0xD0, 0xD1, 0x02, 0x9E, 0x7E, 0x57,
		0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
		0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
	}
)
