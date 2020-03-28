package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
)

func hybridEncryption(text []byte, pubkey *rsa.PublicKey) []byte {
	symmetricKey := generateSymmetricKey(AesKeyBytes) // TODO get from utils
	ciphertext, err := symmetricEncryption(text, symmetricKey)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	cipherSymmetricKey := EncryptWithPublicKey(symmetricKey, pubkey)
	return append(cipherSymmetricKey, ciphertext...)
}

func hybridDecryption(cipherTextAndKey []byte, privkey *rsa.PrivateKey) []byte {
	// TODO check if len less then len
	cipherSymmetricKey := cipherTextAndKey[:CipherRsaLen]
	ciphertext := cipherTextAndKey[CipherRsaLen:]
	plainSymmetricKey := DecryptWithPrivateKey(cipherSymmetricKey, privkey)
	plaintext, err := symmetricDecryption(ciphertext, plainSymmetricKey)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	return plaintext
}

func createKeys(myPublicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, pubkey := GenerateRsaKeyPair(RsaKeyBits)
	WritePublicKeyToFile(myPublicKeyPath, pubkey)
	return privkey, pubkey
}

func generateSymmetricKey(bytes int) []byte {
	rng := rand.Reader
	key := make([]byte, bytes)
	if _, err := io.ReadFull(rng, key); err != nil {
		panic("RNG failure")
	}
	return key
}

func symmetricEncryption(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func symmetricDecryption(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func GenerateRsaKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Fatal(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("not ok")
	}
	return key
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}

func WritePublicKeyToFile(filePath string, key *rsa.PublicKey) {
	writeToFile(filePath, PublicKeyToBytes(key))
}

func WritePrivateKeyToFile(filePath string, key *rsa.PrivateKey) {
	writeToFile(filePath, PrivateKeyToBytes(key))
}

func ReadPublicKeyFromFile(filePath string) *rsa.PublicKey {
	return BytesToPublicKey(readFromFile(filePath))
}

func ReadPrivateKeyFromFile(filePath string) *rsa.PrivateKey {
	return BytesToPrivateKey(readFromFile(filePath))
}
