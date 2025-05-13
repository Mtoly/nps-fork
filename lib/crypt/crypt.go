package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/crypto/blake2b"
)

// en
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// de
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	err, origData = PKCS5UnPadding(origData)
	return origData, err
}

// Completion when the length is insufficient
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Remove excess
func PKCS5UnPadding(origData []byte) (error, []byte) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if (length - unpadding) < 0 {
		return errors.New("len error"), nil
	}
	return nil, origData[:(length - unpadding)]
}

// EncryptBytes AES-GCM
func EncryptBytes(data []byte, keyStr string) ([]byte, error) {
	if keyStr == "" {
		return data, nil
	}
	key := sha256.Sum256([]byte(keyStr))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("io.ReadFull: %w", err)
	}
	ct := gcm.Seal(nil, nonce, data, nil)
	return append(nonce, ct...), nil
}

// DecryptBytes AES-GCM
func DecryptBytes(enc []byte, keyStr string) ([]byte, error) {
	if keyStr == "" {
		return enc, nil
	}
	key := sha256.Sum256([]byte(keyStr))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}
	ns := gcm.NonceSize()
	if len(enc) < ns+gcm.Overhead() {
		return nil, fmt.Errorf("ciphertext too short: %d", len(enc))
	}
	nonce, ct := enc[:ns], enc[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm.Open: %w", err)
	}
	return pt, nil
}

// Get HMAC value
func ComputeHMAC(vkey string, timestamp int64, randomDataPieces ...[]byte) []byte {
	key := []byte(vkey)
	tsBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBuf, uint64(timestamp))
	mac := hmac.New(sha256.New, key)
	mac.Write(tsBuf)
	for _, data := range randomDataPieces {
		mac.Write(data)
	}
	return mac.Sum(nil) // 32bit
}

// Generate 32-bit MD5 strings
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Generate 64-bit BLAKE2b-256 strings
func Blake2b(s string) string {
	hash := blake2b.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// GetRandomString 生成指定长度的随机密钥，支持可选传入id
func GetRandomString(l int, id ...int) string {
	// 字符集
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	dictBytes := []byte(str)
	var result []byte

	// 如果传入id，则将id转换为字符集映射并倒序放在最前面
	if len(id) > 0 {
		// 将id转为字符集表示的字符串
		idMapped := ""
		for id[0] > 0 {
			idMapped = string(str[id[0]%len(str)]) + idMapped
			id[0] /= len(str)
		}

		// 如果倒序后的id长度超过指定长度l，则截断
		//if len(idMapped) > l {
		//	idMapped = idMapped[:l]
		//}

		// 将倒序后的id添加到结果中
		result = append(result, []byte(idMapped)...)
	}

	// 计算剩余需要生成的随机字符的长度
	remainingLength := l - len(result)
	if remainingLength > 0 {
		// 使用当前时间的UnixNano作为随机数种子
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// 生成剩余的随机字符
		for i := 0; i < remainingLength; i++ {
			//result = append(result, dictBytes[r.Intn(len(dictBytes))])
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(dictBytes))))
			if err != nil {
				// 如果安全随机生成失败，回退到时间戳伪随机
				idx := int(time.Now().UnixNano() % int64(len(dictBytes)))
				result = append(result, dictBytes[idx])
				continue
			}
			result = append(result, dictBytes[int(nBig.Int64())])
		}
	}

	// 返回最终结果字符串
	return string(result)
}
