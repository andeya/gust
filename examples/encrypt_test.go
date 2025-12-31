package examples_test

import (
	"fmt"

	"github.com/andeya/gust/encrypt"
)

// Example_encrypt_hashFunctions demonstrates various hash functions.
func Example_encrypt_hashFunctions() {
	data := []byte("hello world")

	// MD5 hash (hex encoding)
	md5Hash := encrypt.MD5(data, encrypt.EncodingHex)
	if md5Hash.IsOk() {
		fmt.Println("MD5 (hex):", md5Hash.Unwrap())
	}

	// SHA256 hash (hex encoding)
	sha256Hash := encrypt.SHA256(data, encrypt.EncodingHex)
	if sha256Hash.IsOk() {
		fmt.Println("SHA256 (hex):", sha256Hash.Unwrap())
	}

	// SHA512 hash (base64 encoding)
	sha512Hash := encrypt.SHA512(data, encrypt.EncodingBase64)
	if sha512Hash.IsOk() {
		fmt.Println("SHA512 (base64) length:", len(sha512Hash.Unwrap()))
	}

	// Output:
	// MD5 (hex): 5eb63bbbe01eeed093cb22bb8f5acdc3
	// SHA256 (hex): b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
	// SHA512 (base64) length: 86
}

// Example_encrypt_aesEncryption demonstrates AES encryption and decryption.
func Example_encrypt_aesEncryption() {
	// AES-128 key (16 bytes)
	key := []byte("1234567890123456")
	plaintext := []byte("secret message")

	// Encrypt with CBC mode and hex encoding
	encrypted := encrypt.EncryptAES(key, plaintext, encrypt.ModeCBC, encrypt.EncodingHex)
	if encrypted.IsErr() {
		fmt.Println("Encryption error:", encrypted.Err())
		return
	}

	ciphertext := encrypted.Unwrap()
	fmt.Println("Encrypted:", len(ciphertext) > 0)

	// Decrypt
	decrypted := encrypt.DecryptAES(key, ciphertext, encrypt.ModeCBC, encrypt.EncodingHex)
	if decrypted.IsOk() {
		fmt.Println("Decrypted:", string(decrypted.Unwrap()))
	}

	// Output:
	// Encrypted: true
	// Decrypted: secret message
}

// Example_encrypt_encodingFormats demonstrates different encoding formats.
func Example_encrypt_encodingFormats() {
	data := []byte("test data")

	// Hex encoding (default, most common)
	hexHash := encrypt.SHA256(data, encrypt.EncodingHex)
	if hexHash.IsOk() {
		fmt.Println("Hex encoding length:", len(hexHash.Unwrap()))
	}

	// Base64 encoding (more compact)
	base64Hash := encrypt.SHA256(data, encrypt.EncodingBase64)
	if base64Hash.IsOk() {
		fmt.Println("Base64 encoding length:", len(base64Hash.Unwrap()))
	}

	// Output:
	// Hex encoding length: 64
	// Base64 encoding length: 43
}

// Example_encrypt_shaVariants demonstrates SHA family variants.
func Example_encrypt_shaVariants() {
	data := []byte("hello")

	// SHA-1
	sha1 := encrypt.SHA1(data, encrypt.EncodingHex)
	if sha1.IsOk() {
		fmt.Println("SHA-1 computed:", len(sha1.Unwrap()) > 0)
	}

	// SHA-224
	sha224 := encrypt.SHA224(data, encrypt.EncodingHex)
	if sha224.IsOk() {
		hash := sha224.Unwrap()
		fmt.Println("SHA-224 computed:", len(hash) > 0)
	}

	// SHA-256
	sha256 := encrypt.SHA256(data, encrypt.EncodingHex)
	if sha256.IsOk() {
		fmt.Println("SHA-256 computed:", len(sha256.Unwrap()) > 0)
	}

	// SHA-384
	sha384 := encrypt.SHA384(data, encrypt.EncodingHex)
	if sha384.IsOk() {
		fmt.Println("SHA-384 computed:", len(sha384.Unwrap()) > 0)
	}

	// SHA-512
	sha512 := encrypt.SHA512(data, encrypt.EncodingHex)
	if sha512.IsOk() {
		fmt.Println("SHA-512 computed:", len(sha512.Unwrap()) > 0)
	}

	// Output:
	// SHA-1 computed: true
	// SHA-224 computed: true
	// SHA-256 computed: true
	// SHA-384 computed: true
	// SHA-512 computed: true
}

// Example_encrypt_fnvHash demonstrates FNV hash functions.
func Example_encrypt_fnvHash() {
	data := []byte("hello")

	// FNV-1a 32-bit
	fnv32 := encrypt.FNV1a32(data)
	if fnv32.IsOk() {
		fmt.Println("FNV-1a 32-bit computed:", fnv32.Unwrap() > 0)
	}

	// FNV-1a 64-bit
	fnv64 := encrypt.FNV1a64(data)
	if fnv64.IsOk() {
		fmt.Println("FNV-1a 64-bit computed:", fnv64.Unwrap() > 0)
	}

	// FNV-1a 128-bit (with encoding)
	fnv128 := encrypt.FNV1a128(data, encrypt.EncodingHex)
	if fnv128.IsOk() {
		fmt.Println("FNV-1a 128-bit (hex) length:", len(fnv128.Unwrap()))
	}

	// Output:
	// FNV-1a 32-bit computed: true
	// FNV-1a 64-bit computed: true
	// FNV-1a 128-bit (hex) length: 32
}

// Example_encrypt_crcHash demonstrates CRC hash functions.
func Example_encrypt_crcHash() {
	data := []byte("hello world")

	// CRC-32
	crc32 := encrypt.CRC32(data)
	if crc32.IsOk() {
		fmt.Println("CRC-32 computed:", crc32.Unwrap() > 0)
	}

	// CRC-64
	crc64 := encrypt.CRC64(data)
	if crc64.IsOk() {
		fmt.Println("CRC-64 computed:", crc64.Unwrap() > 0)
	}

	// Adler-32
	adler32 := encrypt.Adler32(data)
	if adler32.IsOk() {
		fmt.Println("Adler-32 computed:", adler32.Unwrap() > 0)
	}

	// Output:
	// CRC-32 computed: true
	// CRC-64 computed: true
	// Adler-32 computed: true
}

// Example_encrypt_aesModes demonstrates different AES encryption modes.
func Example_encrypt_aesModes() {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	plaintext := []byte("test message")

	// ECB mode (not recommended for most use cases)
	ecbResult := encrypt.EncryptAES(key, plaintext, encrypt.ModeECB, encrypt.EncodingHex)
	if ecbResult.IsOk() {
		fmt.Println("ECB mode: encrypted")
	}

	// CBC mode (recommended for most use cases)
	cbcResult := encrypt.EncryptAES(key, plaintext, encrypt.ModeCBC, encrypt.EncodingHex)
	if cbcResult.IsOk() {
		fmt.Println("CBC mode: encrypted")
	}

	// CTR mode (suitable for streaming)
	ctrResult := encrypt.EncryptAES(key, plaintext, encrypt.ModeCTR, encrypt.EncodingHex)
	if ctrResult.IsOk() {
		fmt.Println("CTR mode: encrypted")
	}

	// Output:
	// ECB mode: encrypted
	// CBC mode: encrypted
	// CTR mode: encrypted
}
