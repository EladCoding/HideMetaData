package mixnet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"io"
)


func GenerateAsymmetricKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	scripts.CheckErrToLog(err)
	pubKey := privKey.PublicKey
	return privKey, pubKey
}


func EncryptKeyForKeyExchange(destPubKey ecdsa.PublicKey) ([]byte, ecdsa.PublicKey) {
	privKey, pubKey := GenerateAsymmetricKeyPair()
	sharedSecret, _ := destPubKey.Curve.ScalarMult(destPubKey.X, destPubKey.Y, privKey.D.Bytes())
	return sharedSecret.Bytes(), pubKey
}


func DecryptKeyForKeyExchange(sourcePubKey ecdsa.PublicKey, privKey *ecdsa.PrivateKey) []byte {
	sharedSecret, _ := sourcePubKey.Curve.ScalarMult(sourcePubKey.X, sourcePubKey.Y, privKey.D.Bytes())
	return sharedSecret.Bytes()
}


func hybridEncription(plaintext []byte, destName string) ([]byte, ecdsa.PublicKey, scripts.SecretKey) {
	destPubKey := scripts.DecodePublicKey(userPubKeyMap[destName])
	sharedSecret, pubKey := EncryptKeyForKeyExchange(*destPubKey)
	cipherMsgData, err := symmetricEncryption(plaintext, sharedSecret)
	scripts.CheckErrToLog(err)
	return cipherMsgData, pubKey, sharedSecret
}


func symmetricEncryption(plaintext []byte, key []byte) ([]byte, error) {
	// use the key to create AES cipher
	c, err := aes.NewCipher(key)
	scripts.CheckErrToLog(err)
	gcm, err := cipher.NewGCM(c)
	scripts.CheckErrToLog(err)
	// create nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	scripts.CheckErrToLog(err)
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}


func symmetricDecryption(ciphertext []byte, key []byte) ([]byte, error) {
	// use the key to create AES cipher
	c, err := aes.NewCipher(key)
	scripts.CheckErrToLog(err)
	gcm, err := cipher.NewGCM(c)
	scripts.CheckErrToLog(err)
	// find nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}


func pkcs7strip(msg []byte, blockSize int) ([]byte, error) {
	msgLen := len(msg)
	if msgLen == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	} else if msgLen%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	} else if msgLen > blockSize {
		return nil, fmt.Errorf("pkcs7: Invalid message size %d", msgLen) // TODO maybe remove this part (if I want to split msg)
	}
	padLen := int(msg[msgLen-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(msg, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return msg[:msgLen-padLen], nil
}


func pkcs7padding(msg []byte, blockSize int) ([]byte, error) {
	if blockSize < 0 || blockSize > 256 {
		return nil, fmt.Errorf("pkcs7: Invalid block size %d", blockSize)
	} else if msgLen:= len(msg); msgLen >= blockSize { // TODO maybe remove this part (if I want to split msg)
		return nil, fmt.Errorf("pkcs7: Invalid message size %d", msgLen)
	} else {
		padLen := blockSize - (msgLen % blockSize)
		padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(msg, padding...), nil
	}
}


func generateRandomBytes(numOfBytes int) []byte {
	randomBytesDate := make([]byte, numOfBytes)
	rand.Read(randomBytesDate)
	return randomBytesDate
}
