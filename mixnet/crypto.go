package mixnet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	//"go.dedis.ch/kyber/group/edwards25519"
	//"go.dedis.ch/kyber/util/random"
	//"go.dedis.ch/kyber/encrypt/ecies"
)


func GenerateAsymmetricKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	CheckErrToLog(err)
	pubKey := privKey.PublicKey
	return privKey, pubKey
}


func EncryptKeyForKeyExchange(destPubKey ecdsa.PublicKey) ([]byte, ecdsa.PublicKey) {
	privKey, pubKey := GenerateAsymmetricKeyPair()
	a, _ := destPubKey.Curve.ScalarMult(destPubKey.X, destPubKey.Y, privKey.D.Bytes())
	sharedSecret := sha256.Sum256(a.Bytes())
	return sharedSecret[:], pubKey
}


func DecryptKeyForKeyExchange(sourcePubKey ecdsa.PublicKey, privKey *ecdsa.PrivateKey) []byte {
	b, _ := sourcePubKey.Curve.ScalarMult(sourcePubKey.X, sourcePubKey.Y, privKey.D.Bytes())
	sharedSecret := sha256.Sum256(b.Bytes())
	return sharedSecret[:]
}


//func hybridEncription(plaintext []byte, destName string) ([]byte) {
//	destPubKey := UserPubKeyMap[destName]
//
//	suite := edwards25519.NewBlakeSHA256Ed25519()
//	cipherMsgData, _ := ecies.Encrypt(suite, destPubKey, plaintext, suite.Hash)
//
//	return cipherMsgData
//}


func hybridEncription(plaintext []byte, destName string) ([]byte, ecdsa.PublicKey, SecretKey) {
	destPubKey := UserPubKeyMap[destName]
	sharedSecret, pubKey := EncryptKeyForKeyExchange(*destPubKey)
	if len(sharedSecret) != 32 {
		print(len(sharedSecret))
		panic("------------------------------------------------------------------------ what the len!!!!!!!!!!!!!!!! ------------------------------------------------------------------------")
	}
	cipherMsgData, err := symmetricEncryption(plaintext, sharedSecret)
	CheckErrToLog(err)
	return cipherMsgData, pubKey, sharedSecret
}


const NONCE_LEN int = 24
func randomNonce() ([]byte, error) {
	b := make([]byte, NONCE_LEN)
	_, err := rand.Read(b)
	return b, err
}


//func symmetricEncryption(plaintext, k []byte) ([]byte, error) {
//	nonce, err := randomNonce()
//	if err != nil {
//		return nil, err
//	}
//	cyphertext := sodium.Bytes(plaintext).SecretBox(
//		sodium.SecretBoxNonce{nonce},
//		sodium.SecretBoxKey{k})
//	return append(nonce, cyphertext...), nil
//}


//func symmetricDecryption(cyphertext, k []byte) ([]byte, error) {
//	nonce := sodium.SecretBoxNonce{cyphertext[:NONCE_LEN]}
//	enc := sodium.Bytes(cyphertext[NONCE_LEN:])
//	return enc.SecretBoxOpen(nonce, sodium.SecretBoxKey{k})
//}


func symmetricEncryption(plaintext []byte, key []byte) ([]byte, error) {
	// use the key to create AES cipher
	c, err := aes.NewCipher(key)
	CheckErrToLog(err)
	gcm, err := cipher.NewGCM(c)
	CheckErrToLog(err)
	// create nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	CheckErrToLog(err)
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}


func symmetricDecryption(ciphertext []byte, key []byte) ([]byte, error) {
	// use the key to create AES cipher
	c, err := aes.NewCipher(key)
	CheckErrToLog(err)
	gcm, err := cipher.NewGCM(c)
	CheckErrToLog(err)
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
		return nil, fmt.Errorf("pkcs7: Invalid message size %d", msgLen)
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
	} else if msgLen:= len(msg); msgLen >= blockSize {
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
