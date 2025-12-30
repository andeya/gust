// Package encrypt provides a comprehensive collection of cryptographic functions
// with Rust-inspired error handling.
//
// This package offers:
//   - Complete hash function suite: MD5, SHA-1, SHA-224, SHA-256, SHA-384, SHA-512,
//     SHA-512/224, SHA-512/256, FNV-1, FNV-1a (32/64/128-bit), CRC-32, CRC-64, Adler-32
//   - AES encryption/decryption with multiple modes (ECB, CBC, CTR)
//   - Support for both hex and base64 encoding formats
//   - Type-safe error handling using gust's Result[T] type for chainable operations
//
// # Examples
//
//	// Hash operations (user can specify encoding)
//	hash := encrypt.SHA256([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
//
//	// AES encryption with chainable operations
//	key := []byte("1234567890123456") // 16 bytes for AES-128
//	plaintext := []byte("secret data")
//
//	result := encrypt.EncryptAES(key, plaintext, encrypt.ModeCBC, encrypt.EncodingBase64).
//		Map(func(ciphertext []byte) []byte {
//			// Further processing if needed
//			return ciphertext
//		})
//
//	if result.IsOk() {
//		ciphertext := result.Unwrap()
//		// Decrypt
//		decrypted := encrypt.DecryptAES(key, ciphertext, encrypt.ModeCBC, encrypt.EncodingBase64)
//		if decrypted.IsOk() {
//			fmt.Println(string(decrypted.Unwrap()))
//		}
//	}
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"

	"github.com/andeya/gust/result"
)

const (
	// aesKeySize128 is the key size for AES-128 in bytes.
	aesKeySize128 = 16
	// aesKeySize192 is the key size for AES-192 in bytes.
	aesKeySize192 = 24
	// aesKeySize256 is the key size for AES-256 in bytes.
	aesKeySize256 = 32
)

// Encoding represents the encoding format for encrypted data and hash output.
type Encoding int

const (
	// EncodingHex uses hexadecimal encoding (default, most common).
	// This is the standard format used by command-line tools like md5sum, sha256sum.
	EncodingHex Encoding = iota
	// EncodingBase64 uses Base64URL encoding (more compact, URL-safe).
	EncodingBase64
)

// Mode represents the AES encryption mode.
type Mode int

const (
	// ModeECB uses Electronic Codebook mode (not recommended for most use cases).
	ModeECB Mode = iota
	// ModeCBC uses Cipher Block Chaining mode (recommended for most use cases).
	ModeCBC
	// ModeCTR uses Counter mode (suitable for streaming).
	ModeCTR
)

// MD5 returns the MD5 checksum of the data as an encoded string.
//
// # Examples
//
//	// Using default hex encoding (most common)
//	hash := encrypt.MD5([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
//
//	// Using base64 encoding (more compact)
//	hash := encrypt.MD5([]byte("hello"), encrypt.EncodingBase64)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // base64-encoded hash string
//	}
func MD5(data []byte, encoding Encoding) result.Result[string] {
	checksum := md5.Sum(data)
	return encodeHash(checksum[:], encoding)
}

// SHA1 returns the SHA1 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA1([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA1(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha1.Sum(data)
	return encodeHash(checksum[:], encoding)
}

// SHA256 returns the SHA256 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA256([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA256(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha256.Sum256(data)
	return encodeHash(checksum[:], encoding)
}

// SHA512 returns the SHA512 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA512([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA512(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha512.Sum512(data)
	return encodeHash(checksum[:], encoding)
}

// SHA224 returns the SHA224 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA224([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA224(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha256.Sum224(data)
	return encodeHash(checksum[:], encoding)
}

// SHA384 returns the SHA384 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA384([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA384(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha512.Sum384(data)
	return encodeHash(checksum[:], encoding)
}

// SHA512_224 returns the SHA-512/224 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA512_224([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA512_224(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha512.Sum512_224(data)
	return encodeHash(checksum[:], encoding)
}

// SHA512_256 returns the SHA-512/256 checksum of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.SHA512_256([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func SHA512_256(data []byte, encoding Encoding) result.Result[string] {
	checksum := sha512.Sum512_256(data)
	return encodeHash(checksum[:], encoding)
}

// FNV1a64 returns the 64-bit FNV-1a hash of the data.
//
// # Examples
//
//	hash := encrypt.FNV1a64([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint64 value
//	}
func FNV1a64(data []byte) result.Result[uint64] {
	h := fnv.New64a()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[uint64](err)
	}
	return result.Ok(h.Sum64())
}

// FNV1a32 returns the 32-bit FNV-1a hash of the data.
//
// # Examples
//
//	hash := encrypt.FNV1a32([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint32 value
//	}
func FNV1a32(data []byte) result.Result[uint32] {
	h := fnv.New32a()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[uint32](err)
	}
	return result.Ok(h.Sum32())
}

// FNV1_32 returns the 32-bit FNV-1 hash of the data.
//
// # Examples
//
//	hash := encrypt.FNV1_32([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint32 value
//	}
func FNV1_32(data []byte) result.Result[uint32] {
	h := fnv.New32()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[uint32](err)
	}
	return result.Ok(h.Sum32())
}

// FNV1_64 returns the 64-bit FNV-1 hash of the data.
//
// # Examples
//
//	hash := encrypt.FNV1_64([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint64 value
//	}
func FNV1_64(data []byte) result.Result[uint64] {
	h := fnv.New64()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[uint64](err)
	}
	return result.Ok(h.Sum64())
}

// FNV1a128 returns the 128-bit FNV-1a hash of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.FNV1a128([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func FNV1a128(data []byte, encoding Encoding) result.Result[string] {
	h := fnv.New128a()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[string](err)
	}
	return encodeHash(h.Sum(nil), encoding)
}

// FNV1_128 returns the 128-bit FNV-1 hash of the data as an encoded string.
//
// # Examples
//
//	hash := encrypt.FNV1_128([]byte("hello"), encrypt.EncodingHex)
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // hex-encoded hash string
//	}
func FNV1_128(data []byte, encoding Encoding) result.Result[string] {
	h := fnv.New128()
	_, err := h.Write(data)
	if err != nil {
		return result.TryErr[string](err)
	}
	return encodeHash(h.Sum(nil), encoding)
}

// CRC32 returns the CRC-32 checksum of the data.
//
// # Examples
//
//	hash := encrypt.CRC32([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint32 value
//	}
func CRC32(data []byte) result.Result[uint32] {
	return result.Ok(crc32.ChecksumIEEE(data))
}

// CRC64 returns the CRC-64 checksum of the data.
//
// # Examples
//
//	hash := encrypt.CRC64([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint64 value
//	}
func CRC64(data []byte) result.Result[uint64] {
	return result.Ok(crc64.Checksum(data, crc64.MakeTable(crc64.ISO)))
}

// Adler32 returns the Adler-32 checksum of the data.
//
// # Examples
//
//	hash := encrypt.Adler32([]byte("hello"))
//	if hash.IsOk() {
//		fmt.Println(hash.Unwrap()) // uint32 value
//	}
func Adler32(data []byte) result.Result[uint32] {
	return result.Ok(adler32.Checksum(data))
}

// EncryptAES encrypts data using AES with the specified mode and encoding.
//
// The key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 respectively.
//
// # Examples
//
//	key := []byte("1234567890123456") // 16 bytes for AES-128
//	plaintext := []byte("secret data")
//
//	result := encrypt.EncryptAES(key, plaintext, encrypt.ModeCBC, encrypt.EncodingBase64)
//	if result.IsOk() {
//		ciphertext := result.Unwrap()
//		// Use ciphertext...
//	}
func EncryptAES(key, plaintext []byte, mode Mode, encoding Encoding) result.Result[[]byte] {
	if err := validateKey(key); err != nil {
		return result.TryErr[[]byte](err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return result.TryErr[[]byte](err)
	}

	var ciphertextResult result.Result[[]byte]
	switch mode {
	case ModeECB:
		ciphertextResult = encryptECB(block, plaintext)
	case ModeCBC:
		ciphertextResult = encryptCBC(block, plaintext)
	case ModeCTR:
		ciphertextResult = encryptCTR(block, plaintext)
	default:
		return result.FmtErr[[]byte]("unsupported encryption mode: %d", mode)
	}

	if ciphertextResult.IsErr() {
		return ciphertextResult
	}

	return encodeBytes(ciphertextResult.Unwrap(), encoding)
}

// DecryptAES decrypts data using AES with the specified mode and encoding.
//
// The key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 respectively.
// The encoding must match the encoding used during encryption.
//
// # Examples
//
//	key := []byte("1234567890123456")
//	ciphertext := []byte("encrypted data")
//
//	result := encrypt.DecryptAES(key, ciphertext, encrypt.ModeCBC, encrypt.EncodingBase64)
//	if result.IsOk() {
//		plaintext := result.Unwrap()
//		fmt.Println(string(plaintext))
//	}
func DecryptAES(key, ciphertext []byte, mode Mode, encoding Encoding) result.Result[[]byte] {
	if err := validateKey(key); err != nil {
		return result.TryErr[[]byte](err)
	}

	decodedResult := decodeBytes(ciphertext, encoding)
	if decodedResult.IsErr() {
		return decodedResult
	}
	decoded := decodedResult.Unwrap()

	block, err := aes.NewCipher(key)
	if err != nil {
		return result.TryErr[[]byte](err)
	}

	var plaintextResult result.Result[[]byte]
	switch mode {
	case ModeECB:
		plaintextResult = decryptECB(block, decoded)
	case ModeCBC:
		plaintextResult = decryptCBC(block, decoded)
	case ModeCTR:
		plaintextResult = decryptCTR(block, decoded)
	default:
		return result.FmtErr[[]byte]("unsupported decryption mode: %d", mode)
	}

	return plaintextResult
}

// validateKey validates that the key length is valid for AES.
func validateKey(key []byte) error {
	keyLen := len(key)
	if keyLen != aesKeySize128 && keyLen != aesKeySize192 && keyLen != aesKeySize256 {
		return errors.New("AES key must be 16, 24, or 32 bytes")
	}
	return nil
}

// encryptECB encrypts data using ECB mode.
func encryptECB(block cipher.Block, plaintext []byte) result.Result[[]byte] {
	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	ciphertext := make([]byte, len(plaintext))
	dst := ciphertext
	for len(plaintext) > 0 {
		block.Encrypt(dst, plaintext[:blockSize])
		plaintext = plaintext[blockSize:]
		dst = dst[blockSize:]
	}
	return result.Ok(ciphertext)
}

// decryptECB decrypts data using ECB mode.
func decryptECB(block cipher.Block, ciphertext []byte) result.Result[[]byte] {
	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return result.TryErr[[]byte]("ciphertext is not a multiple of the block size")
	}
	plaintext := make([]byte, len(ciphertext))
	dst := plaintext
	for len(ciphertext) > 0 {
		block.Decrypt(dst, ciphertext[:blockSize])
		ciphertext = ciphertext[blockSize:]
		dst = dst[blockSize:]
	}
	return pkcs5Unpadding(plaintext)
}

// encryptCBC encrypts data using CBC mode.
func encryptCBC(block cipher.Block, plaintext []byte) result.Result[[]byte] {
	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return result.TryErr[[]byte](err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return result.Ok(ciphertext)
}

// decryptCBC decrypts data using CBC mode.
func decryptCBC(block cipher.Block, ciphertext []byte) result.Result[[]byte] {
	if len(ciphertext) < aes.BlockSize {
		return result.TryErr[[]byte]("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return result.TryErr[[]byte]("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	plaintext := ciphertext
	mode.CryptBlocks(plaintext, ciphertext)
	return pkcs5Unpadding(plaintext)
}

// encryptCTR encrypts data using CTR mode.
func encryptCTR(block cipher.Block, plaintext []byte) result.Result[[]byte] {
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return result.TryErr[[]byte](err)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return result.Ok(ciphertext)
}

// decryptCTR decrypts data using CTR mode.
func decryptCTR(block cipher.Block, ciphertext []byte) result.Result[[]byte] {
	if len(ciphertext) < aes.BlockSize {
		return result.TryErr[[]byte]("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	plaintext := ciphertext
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(plaintext, ciphertext)
	return result.Ok(plaintext)
}

// pkcs5Padding adds PKCS5 padding to the plaintext.
func pkcs5Padding(plaintext []byte, blockSize int) []byte {
	n := byte(blockSize - len(plaintext)%blockSize)
	for i := byte(0); i < n; i++ {
		plaintext = append(plaintext, n)
	}
	return plaintext
}

// pkcs5Unpadding removes PKCS5 padding from the plaintext.
func pkcs5Unpadding(data []byte) result.Result[[]byte] {
	l := len(data)
	if l == 0 {
		return result.TryErr[[]byte]("input padded bytes is empty")
	}
	padLen := int(data[l-1])
	if padLen == 0 || padLen > l {
		return result.TryErr[[]byte]("invalid padding length")
	}
	if l-padLen < 0 {
		return result.TryErr[[]byte]("input padded bytes is invalid")
	}
	pad := data[l-padLen : l]
	for _, v := range pad {
		if v != byte(padLen) {
			return result.TryErr[[]byte]("invalid padding bytes")
		}
	}
	return result.Ok(data[:l-padLen])
}

// encodeBytes encodes data using the specified encoding.
func encodeBytes(data []byte, encoding Encoding) result.Result[[]byte] {
	switch encoding {
	case EncodingHex:
		dst := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(dst, data)
		return result.Ok(dst)
	case EncodingBase64:
		buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(data)))
		base64.RawURLEncoding.Encode(buf, data)
		return result.Ok(buf)
	default:
		return result.FmtErr[[]byte]("unsupported encoding: %d", encoding)
	}
}

// decodeBytes decodes data using the specified encoding.
func decodeBytes(data []byte, encoding Encoding) result.Result[[]byte] {
	switch encoding {
	case EncodingHex:
		dst := make([]byte, hex.DecodedLen(len(data)))
		n, err := hex.Decode(dst, data)
		if err != nil {
			return result.TryErr[[]byte](err)
		}
		return result.Ok(dst[:n])
	case EncodingBase64:
		dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
		n, err := base64.RawURLEncoding.Decode(dst, data)
		if err != nil {
			return result.TryErr[[]byte](err)
		}
		return result.Ok(dst[:n])
	default:
		return result.FmtErr[[]byte]("unsupported encoding: %d", encoding)
	}
}

// encodeHash encodes hash data using the specified encoding format.
func encodeHash(data []byte, encoding Encoding) result.Result[string] {
	switch encoding {
	case EncodingHex:
		return result.Ok(hex.EncodeToString(data))
	case EncodingBase64:
		return result.Ok(base64.RawURLEncoding.EncodeToString(data))
	default:
		return result.FmtErr[string]("unsupported encoding: %d", encoding)
	}
}
