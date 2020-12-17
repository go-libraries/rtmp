package handshake

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

func HMACSha256(msgBytes []byte, key []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(msgBytes)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Get timestamp  4 byte
func GetTimestampByte() []byte {
	//return int(0)
	bt := make([]byte, 8)
	binary.BigEndian.PutUint64(bt, uint64(time.Now().UnixNano()/int64(1000000))%MaxTimeStamp)
	return bt[4:]
}

func validateDigest(buf []byte, offset int, key []byte) int {
	digestPos := calOffset(buf, offset, offset+4)
	// Create temp buffer
	tmpBuf := new(bytes.Buffer)
	tmpBuf.Write(buf[:digestPos])
	tmpBuf.Write(buf[digestPos+RtmpSha256DigestLength:])
	// Generate the hash
	tempHash, err := HMACSha256(tmpBuf.Bytes(), key)
	if err != nil {
		return 0
	}
	if bytes.Compare(tempHash, buf[digestPos:digestPos+RtmpSha256DigestLength]) == 0 {
		return digestPos
	}
	return 0
}

func getS1SeverPost(buf []byte) (pos int) {
	//由于
	pos = validateDigest(buf, 8, GenuineFmsKey[:36])
	if pos == 0 {
		// (1536 - 8) / 2 + 8 = 772
		return validateDigest(buf, 772, GenuineFmsKey[:36])
	}

	return 0
}

//digestPos := CalcDigestPos(buf, offset, 728, offset+4)

//func CalcDigestPos(buf []byte, offset int, mod_val int, add_val int) (digest_pos int) {
//	var i int
//	for i = 0; i < 4; i++ {
//		digest_pos += int(buf[i+offset])
//	}
//	digest_pos = digest_pos%mod_val + add_val
//	return
//}

/**
offs = ngx_rtmp_find_digest(b, peer_key, 772, s->connection->log);
   if (offs == NGX_ERROR) {
       offs = ngx_rtmp_find_digest(b, peer_key, 8, s->connection->log);
   }
ngx_rtmp_find_digest(ngx_buf_t *b, ngx_str_t *key, size_t base, ngx_log_t *log)
{
    size_t                  n, offs;
    u_char                  digest[NGX_RTMP_HANDSHAKE_KEYLEN];
    u_char                 *p;

    offs = 0;
    for (n = 0; n < 4; ++n) {
        offs += b->pos[base + n];
    }
    offs = (offs % 728) + base + 4;
    p = b->pos + offs;

    if (ngx_rtmp_make_digest(key, b, p, digest, log) != NGX_OK) {
        return NGX_ERROR;
    }

    if (ngx_memcmp(digest, p, NGX_RTMP_HANDSHAKE_KEYLEN) == 0) {
        return offs;
    }

    return NGX_ERROR;
}
*/

//b 原始b
//key 混淆k
//基础偏移量（兼容2种数据模式）
func buildC1Digest(buf []byte, key []byte) ([]byte, int) {
	//time 4 version 4 offset 4 => addVal 12
	pos := calcDigestPos(buf, 8, 728)
	tmpBuf := new(bytes.Buffer)
	tmpBuf.Write(buf[:pos])
	tmpBuf.Write(buf[pos+RtmpSha256DigestLength:])
	// Generate the hash
	tempHash, _ := HMACSha256(tmpBuf.Bytes(), key)
	copy(buf[pos:], tempHash)

	return buf, pos
}

func calcDigestPos(buf []byte, offset int, modVal int) (digestPos int) {
	var i int
	for i = 0; i < 4; i++ {
		digestPos += int(buf[i+offset])
	}
	digestPos = digestPos%modVal + offset+4
	return
}

//计算公式
func calOffset(b []byte, offset int, addVal int) (pos int) {
	//找到offset那四个字节分别转换成数字相加
	//因为这个是随机数有可能超过了key的所有长度所有需要取余数。
	//这个取余的数就是764字节的整体长度- 128字节的Key秘钥-4字节的自身长度
	// modelVal
	// modelVal = 764 - 128 -4
	var i int
	for ; i < 4; i++ {
		pos += int(b[i+offset])
	}
	pos = pos%728 + addVal
	return
	//for n := 8; n < 12; n++ { //time + version 0->7  offset len(4)
	//	pos += int(b[offset+n])
	//}
	//pos = (pos % 728) + offset //offset = offset sum(bt[8:12]) % 728 + base ,
	//return
}
