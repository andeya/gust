package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := MD5(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := md5.Sum(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := MD5(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)

	// Test empty input
	result3 := MD5([]byte{}, EncodingHex)
	assert.True(t, result3.IsOk())
	assert.NotEmpty(t, result3.Unwrap())
}

func TestSHA1(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA1(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha1.Sum(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA1(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestSHA256(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA256(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha256.Sum256(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA256(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestSHA512(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA512(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha512.Sum512(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA512(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestFNV1a64(t *testing.T) {
	data := []byte("hello")

	result := FNV1a64(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	h := fnv.New64a()
	h.Write(data)
	expected := h.Sum64()
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := FNV1a64([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint64(0xcbf29ce484222325), result2.Unwrap()) // FNV-1a empty hash
}

func TestFNV1a32(t *testing.T) {
	data := []byte("hello")

	result := FNV1a32(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	h := fnv.New32a()
	h.Write(data)
	expected := h.Sum32()
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := FNV1a32([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint32(0x811c9dc5), result2.Unwrap()) // FNV-1a empty hash
}

func TestEncryptAES_ECB(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	plaintext := []byte("hello world")

	// Test hex encoding
	result := EncryptAES(key, plaintext, ModeECB, EncodingHex)
	assert.True(t, result.IsOk())
	ciphertext := result.Unwrap()

	// Decrypt and verify
	decrypted := DecryptAES(key, ciphertext, ModeECB, EncodingHex)
	assert.True(t, decrypted.IsOk())
	assert.Equal(t, plaintext, decrypted.Unwrap())

	// Test base64 encoding
	result2 := EncryptAES(key, plaintext, ModeECB, EncodingBase64)
	assert.True(t, result2.IsOk())
	ciphertext2 := result2.Unwrap()

	// Decrypt and verify
	decrypted2 := DecryptAES(key, ciphertext2, ModeECB, EncodingBase64)
	assert.True(t, decrypted2.IsOk())
	assert.Equal(t, plaintext, decrypted2.Unwrap())
}

func TestEncryptAES_CBC(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	plaintext := []byte("hello world")

	// Test hex encoding
	result := EncryptAES(key, plaintext, ModeCBC, EncodingHex)
	assert.True(t, result.IsOk())
	ciphertext := result.Unwrap()

	// Decrypt and verify
	decrypted := DecryptAES(key, ciphertext, ModeCBC, EncodingHex)
	assert.True(t, decrypted.IsOk())
	assert.Equal(t, plaintext, decrypted.Unwrap())

	// Test base64 encoding
	result2 := EncryptAES(key, plaintext, ModeCBC, EncodingBase64)
	assert.True(t, result2.IsOk())
	ciphertext2 := result2.Unwrap()

	// Decrypt and verify
	decrypted2 := DecryptAES(key, ciphertext2, ModeCBC, EncodingBase64)
	assert.True(t, decrypted2.IsOk())
	assert.Equal(t, plaintext, decrypted2.Unwrap())

	// Test that same plaintext produces different ciphertext (due to random IV)
	result3 := EncryptAES(key, plaintext, ModeCBC, EncodingHex)
	assert.True(t, result3.IsOk())
	ciphertext3 := result3.Unwrap()
	// Should be different due to random IV
	assert.NotEqual(t, ciphertext, ciphertext3)
}

func TestEncryptAES_CTR(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	plaintext := []byte("hello world")

	// Test hex encoding
	result := EncryptAES(key, plaintext, ModeCTR, EncodingHex)
	assert.True(t, result.IsOk())
	ciphertext := result.Unwrap()

	// Decrypt and verify
	decrypted := DecryptAES(key, ciphertext, ModeCTR, EncodingHex)
	assert.True(t, decrypted.IsOk())
	assert.Equal(t, plaintext, decrypted.Unwrap())

	// Test base64 encoding
	result2 := EncryptAES(key, plaintext, ModeCTR, EncodingBase64)
	assert.True(t, result2.IsOk())
	ciphertext2 := result2.Unwrap()

	// Decrypt and verify
	decrypted2 := DecryptAES(key, ciphertext2, ModeCTR, EncodingBase64)
	assert.True(t, decrypted2.IsOk())
	assert.Equal(t, plaintext, decrypted2.Unwrap())
}

func TestEncryptAES_KeySizes(t *testing.T) {
	plaintext := []byte("hello world")

	testCases := []struct {
		name string
		key  []byte
	}{
		{"AES-128", []byte("1234567890123456")},                 // 16 bytes
		{"AES-192", []byte("123456789012345678901234")},         // 24 bytes
		{"AES-256", []byte("12345678901234567890123456789012")}, // 32 bytes
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EncryptAES(tc.key, plaintext, ModeCBC, EncodingHex)
			assert.True(t, result.IsOk())

			decrypted := DecryptAES(tc.key, result.Unwrap(), ModeCBC, EncodingHex)
			assert.True(t, decrypted.IsOk())
			assert.Equal(t, plaintext, decrypted.Unwrap())
		})
	}
}

func TestEncryptAES_InvalidKey(t *testing.T) {
	plaintext := []byte("hello world")

	invalidKeys := [][]byte{
		[]byte("short"),                             // too short
		[]byte("123456789012345"),                   // 15 bytes
		[]byte("12345678901234567"),                 // 17 bytes
		[]byte("1234567890123456789012345"),         // 25 bytes
		[]byte("123456789012345678901234567890123"), // 33 bytes
	}

	for _, key := range invalidKeys {
		result := EncryptAES(key, plaintext, ModeCBC, EncodingHex)
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "AES key must be 16, 24, or 32 bytes")
	}
}

func TestEncryptAES_EmptyPlaintext(t *testing.T) {
	key := []byte("1234567890123456")

	result := EncryptAES(key, []byte{}, ModeCBC, EncodingHex)
	assert.True(t, result.IsOk())

	decrypted := DecryptAES(key, result.Unwrap(), ModeCBC, EncodingHex)
	assert.True(t, decrypted.IsOk())
	assert.Equal(t, []byte{}, decrypted.Unwrap())
}

func TestDecryptAES_InvalidCiphertext(t *testing.T) {
	key := []byte("1234567890123456")

	// Too short ciphertext
	result := DecryptAES(key, []byte("short"), ModeCBC, EncodingHex)
	assert.True(t, result.IsErr())

	// Invalid hex encoding
	result2 := DecryptAES(key, []byte("invalid hex!"), ModeCBC, EncodingHex)
	assert.True(t, result2.IsErr())

	// Invalid base64 encoding
	result3 := DecryptAES(key, []byte("invalid base64!"), ModeCBC, EncodingBase64)
	assert.True(t, result3.IsErr())
}

func TestDecryptAES_WrongEncoding(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := []byte("hello world")

	// Encrypt with hex
	result := EncryptAES(key, plaintext, ModeCBC, EncodingHex)
	assert.True(t, result.IsOk())
	ciphertext := result.Unwrap()

	// Try to decrypt with base64 (should fail)
	decrypted := DecryptAES(key, ciphertext, ModeCBC, EncodingBase64)
	assert.True(t, decrypted.IsErr())
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
		mode      Mode
		encoding  Encoding
	}{
		{"ECB-Hex", []byte("1234567890123456"), []byte("hello"), ModeECB, EncodingHex},
		{"ECB-Base64", []byte("1234567890123456"), []byte("hello"), ModeECB, EncodingBase64},
		{"CBC-Hex", []byte("1234567890123456"), []byte("hello world"), ModeCBC, EncodingHex},
		{"CBC-Base64", []byte("1234567890123456"), []byte("hello world"), ModeCBC, EncodingBase64},
		{"CTR-Hex", []byte("1234567890123456"), []byte("hello world"), ModeCTR, EncodingHex},
		{"CTR-Base64", []byte("1234567890123456"), []byte("hello world"), ModeCTR, EncodingBase64},
		{"Long text", []byte("1234567890123456"), []byte("This is a longer text to test encryption and decryption"), ModeCBC, EncodingHex},
		{"Empty", []byte("1234567890123456"), []byte{}, ModeCBC, EncodingHex},
		{"AES-192", []byte("123456789012345678901234"), []byte("test"), ModeCBC, EncodingHex},
		{"AES-256", []byte("12345678901234567890123456789012"), []byte("test"), ModeCBC, EncodingHex},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EncryptAES(tc.key, tc.plaintext, tc.mode, tc.encoding)
			assert.True(t, result.IsOk(), "Encryption should succeed for %s", tc.name)

			decrypted := DecryptAES(tc.key, result.Unwrap(), tc.mode, tc.encoding)
			assert.True(t, decrypted.IsOk(), "Decryption should succeed for %s", tc.name)
			assert.Equal(t, tc.plaintext, decrypted.Unwrap(), "Decrypted text should match original for %s", tc.name)
		})
	}
}

func TestHashConsistency(t *testing.T) {
	data := []byte("test data")

	// Same input should produce same hash
	hash1 := MD5(data, EncodingHex)
	hash2 := MD5(data, EncodingHex)
	assert.True(t, hash1.IsOk())
	assert.True(t, hash2.IsOk())
	assert.Equal(t, hash1.Unwrap(), hash2.Unwrap())

	hash3 := SHA256(data, EncodingHex)
	hash4 := SHA256(data, EncodingHex)
	assert.True(t, hash3.IsOk())
	assert.True(t, hash4.IsOk())
	assert.Equal(t, hash3.Unwrap(), hash4.Unwrap())
}

func TestHashDifferentEncodings(t *testing.T) {
	data := []byte("test data")

	// Hex and base64 should produce different strings but same underlying hash
	hashHex := MD5(data, EncodingHex)
	hashBase64 := MD5(data, EncodingBase64)

	assert.True(t, hashHex.IsOk())
	assert.True(t, hashBase64.IsOk())

	// Decode both and verify they represent the same bytes
	hexBytes, err := hex.DecodeString(hashHex.Unwrap())
	assert.NoError(t, err)

	base64Bytes, err := base64.RawURLEncoding.DecodeString(hashBase64.Unwrap())
	assert.NoError(t, err)

	// They should decode to the same bytes
	expected := md5.Sum(data)
	assert.Equal(t, expected[:], hexBytes)
	assert.Equal(t, expected[:], base64Bytes)
}

func TestSHA224(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA224(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha256.Sum224(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA224(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestSHA384(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA384(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha512.Sum384(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA384(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestSHA512_224(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA512_224(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha512.Sum512_224(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA512_224(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestSHA512_256(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := SHA512_256(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)

	// Verify against standard library
	expected := sha512.Sum512_256(data)
	expectedHex := hex.EncodeToString(expected[:])
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := SHA512_256(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected[:])
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestFNV1_32(t *testing.T) {
	data := []byte("hello")

	result := FNV1_32(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	h := fnv.New32()
	h.Write(data)
	expected := h.Sum32()
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := FNV1_32([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint32(0x811c9dc5), result2.Unwrap()) // FNV-1 empty hash
}

func TestFNV1_64(t *testing.T) {
	data := []byte("hello")

	result := FNV1_64(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	h := fnv.New64()
	h.Write(data)
	expected := h.Sum64()
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := FNV1_64([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint64(0xcbf29ce484222325), result2.Unwrap()) // FNV-1 empty hash
}

func TestFNV1a128(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := FNV1a128(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(hashHex)) // 128 bits = 16 bytes = 32 hex chars

	// Verify against standard library
	h := fnv.New128a()
	h.Write(data)
	expected := h.Sum(nil)
	expectedHex := hex.EncodeToString(expected)
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := FNV1a128(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected)
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestFNV1_128(t *testing.T) {
	data := []byte("hello")

	// Test hex encoding
	result := FNV1_128(data, EncodingHex)
	assert.True(t, result.IsOk())
	hashHex := result.Unwrap()

	// Verify it's valid hex
	_, err := hex.DecodeString(hashHex)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(hashHex)) // 128 bits = 16 bytes = 32 hex chars

	// Verify against standard library
	h := fnv.New128()
	h.Write(data)
	expected := h.Sum(nil)
	expectedHex := hex.EncodeToString(expected)
	assert.Equal(t, expectedHex, hashHex)

	// Test base64 encoding
	result2 := FNV1_128(data, EncodingBase64)
	assert.True(t, result2.IsOk())
	hashBase64 := result2.Unwrap()

	// Verify it's valid base64
	_, err = base64.RawURLEncoding.DecodeString(hashBase64)
	assert.NoError(t, err)

	// Verify against standard library
	expectedBase64 := base64.RawURLEncoding.EncodeToString(expected)
	assert.Equal(t, expectedBase64, hashBase64)
}

func TestCRC32(t *testing.T) {
	data := []byte("hello")

	result := CRC32(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	expected := crc32.ChecksumIEEE(data)
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := CRC32([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint32(0), result2.Unwrap())

	// Test consistency
	result3 := CRC32(data)
	result4 := CRC32(data)
	assert.True(t, result3.IsOk())
	assert.True(t, result4.IsOk())
	assert.Equal(t, result3.Unwrap(), result4.Unwrap())
}

func TestCRC64(t *testing.T) {
	data := []byte("hello")

	result := CRC64(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	table := crc64.MakeTable(crc64.ISO)
	expected := crc64.Checksum(data, table)
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := CRC64([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint64(0), result2.Unwrap())

	// Test consistency
	result3 := CRC64(data)
	result4 := CRC64(data)
	assert.True(t, result3.IsOk())
	assert.True(t, result4.IsOk())
	assert.Equal(t, result3.Unwrap(), result4.Unwrap())
}

func TestAdler32(t *testing.T) {
	data := []byte("hello")

	result := Adler32(data)
	assert.True(t, result.IsOk())
	hash := result.Unwrap()

	// Verify against standard library
	expected := adler32.Checksum(data)
	assert.Equal(t, expected, hash)

	// Test empty input
	result2 := Adler32([]byte{})
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint32(1), result2.Unwrap()) // Adler-32 empty hash is 1

	// Test consistency
	result3 := Adler32(data)
	result4 := Adler32(data)
	assert.True(t, result3.IsOk())
	assert.True(t, result4.IsOk())
	assert.Equal(t, result3.Unwrap(), result4.Unwrap())
}
