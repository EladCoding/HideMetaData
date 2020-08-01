package mixnet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	//"go.dedis.ch/kyber/group/edwards25519"
	//"go.dedis.ch/kyber/util/random"
	//"go.dedis.ch/kyber/encrypt/ecies"
)

// Generate ecdsa P256 Assymetric key pair.
func GenerateAsymmetricKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	CheckErrToLog(err)
	pubKey := privKey.PublicKey
	return privKey, pubKey
}

// Encrypt ecdsa P256 public key.
func EncryptKeyForKeyExchange(destPubKey ecdsa.PublicKey) ([]byte, ecdsa.PublicKey) {
	privKey, pubKey := GenerateAsymmetricKeyPair()
	a, _ := destPubKey.Curve.ScalarMult(destPubKey.X, destPubKey.Y, privKey.D.Bytes())
	sharedSecret := sha256.Sum256(a.Bytes())
	return sharedSecret[:], pubKey
}

// Decrypt ecdsa P256 public key.
func DecryptKeyForKeyExchange(sourcePubKey ecdsa.PublicKey, privKey *ecdsa.PrivateKey) []byte {
	b, _ := sourcePubKey.Curve.ScalarMult(sourcePubKey.X, sourcePubKey.Y, privKey.D.Bytes())
	sharedSecret := sha256.Sum256(b.Bytes())
	return sharedSecret[:]
}

// Encrypt a message ecdsa P256 public key, as a hybrid encryption.
func hybridEncryption(plaintext []byte, destName string) ([]byte, ecdsa.PublicKey, SecretKey) {
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

// Encrypt a message, using a key, as a symmetric encryption.
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

// Decrypt a message, using a key, as a symmetric encryption.
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

// Pad bytes to message, using pkcs7 padding.
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

// Strip bytes that was padded by pkcs7 to a message.
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

// Encode ecdsa P256 private key to bytes.
func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}

// Encode ecdsa P256 public key to bytes.
func EncodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return string(pemEncodedPub)
}

// Decode bytes to ecdsa P256 private key.
func DecodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

// Decode bytes to ecdsa P256 public key.
func DecodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}
