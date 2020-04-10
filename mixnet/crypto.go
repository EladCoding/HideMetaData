package mixnet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
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


func hybridEncription(plaintext []byte, destName string) ([]byte, ecdsa.PublicKey) {
	destPubKey := scripts.DecodePublicKey(userPubKeyMap[destName])
	sharedSecret, pubKey := EncryptKeyForKeyExchange(*destPubKey)
	cipherMsgData, err := symmetricEncryption(plaintext, sharedSecret)
	scripts.CheckErrToLog(err)
	return cipherMsgData, pubKey
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
