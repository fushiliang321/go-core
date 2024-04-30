package encrypt

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/savsgio/gotils/strconv"
	"hash"
)

const chunkSizeMax = 512 * 1024 //每个块最大值

func HashChunkWrite(h hash.Hash, bytes []byte) string {
	var (
		bytesLen = len(bytes)
	)
	if bytesLen > chunkSizeMax {
		var (
			start int
			end   int
		)
		for {
			end = start + chunkSizeMax
			if end < bytesLen {
				h.Write(bytes[start:end])
				start = end
			} else {
				h.Write(bytes[start:bytesLen])
				break
			}
		}
	} else {
		h.Write(bytes)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(bytes []byte) string {
	return HashChunkWrite(md5.New(), bytes)
}

// 字符串md5
func StringMD5(v string) string {
	return HashChunkWrite(md5.New(), strconv.S2B(v))
}

func Sha256(bytes []byte) string {
	return HashChunkWrite(sha256.New(), bytes)
}

func StringSha256(v string) string {
	return HashChunkWrite(sha256.New(), strconv.S2B(v))
}
