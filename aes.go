package gcmf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// AesECB 管理对象
var AesECB = new(aesFunc)

type aesFunc struct{}

// eCBEncrypted represents an electronic codebook (ECB) encrypted.
type eCBEncrypted struct {
	block     cipher.Block
	blockSize int
}
type eCBDecrypted struct {
	block     cipher.Block
	blockSize int
}

func (a *aesFunc) Encrypt(input, key []byte) (string, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 padding
	input, err = pkcs7Pad(input, block.BlockSize())
	if err != nil {
		return "", err
	}
	encrypted := make([]byte, len(input))
	ecb := newECBEncrypted(block)
	ecb.CryptBlocks(encrypted, input)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
func (a *aesFunc) Decrypt(ciphertext string, key []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := newECBDecrypted(block)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)
	outByte, err := pkcs7UnPad(ciphertextBytes)
	if err != nil {
		return nil, err
	}
	return outByte, nil
}

// newECBEncrypted creates a new ECBEncrypted.
func newECBEncrypted(b cipher.Block) cipher.BlockMode {
	return &eCBEncrypted{b, b.BlockSize()}
}

// newECBDecrypted returns a new ECB mode decrypted.
func newECBDecrypted(b cipher.Block) cipher.BlockMode {
	return &eCBDecrypted{b, b.BlockSize()}
}

// BlockSize returns the block size of the cipher.
func (x *eCBEncrypted) BlockSize() int { return x.block.BlockSize() }
func (x *eCBDecrypted) BlockSize() int { return x.block.BlockSize() }

// CryptBlocks encrypts or decrypts a number of blocks.
func (x *eCBEncrypted) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("input not full blocks")
	}
	if len(dst) < len(src) {
		panic("output smaller than input")
	}

	for len(src) > 0 {
		x.block.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

}

// CryptBlocks encrypts or decrypts a number of blocks.
func (x *eCBDecrypted) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("input not full blocks")
	}
	if len(dst) < len(src) {
		panic("output smaller than input")
	}

	for len(src) > 0 {
		x.block.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

}
func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: empty input")
	}
	padding := blockSize - (length % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...), nil
}
func pkcs7UnPad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: empty input")
	}
	unPad := int(data[length-1])
	return data[:(length - unPad)], nil
}
