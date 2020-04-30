package mixnet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
)

type UserAddressMapType map[string]string
type UserPublicKeyMapType map[string]string
type UserPrivateKeyMapType map[string]string
type ClientsMap map[string]*rpc.Client
type SecretKey []byte
type EncryptedMsg []byte

var WorkingDir, _ = os.Getwd()
var ExternalsPath = path.Join(WorkingDir, "externals")
var UserAddressesMapPath = path.Join(ExternalsPath, "UserAddressesMap.json")
var UserPublicKeysMapPath = path.Join(ExternalsPath, "userPublicKeysMap.json")
var UserPrivateKeysMapPath = path.Join(ExternalsPath, "userPrivateKeysMap.json")
var PortFormat = "8%s"
var AddressFormat = "localhost:" + PortFormat

var ServerNames = []string{"001", "002", "003"}
var MediatorNames = []string{"101", "102", "103"}
var ClientNames = []string{"201", "202", "203"}
var UserNames = append(append(ServerNames, MediatorNames...), ClientNames...)


func WriteToFile(filePath string, data []byte) {
	// write the whole body at once
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	err := ioutil.WriteFile(filePath, data, 0644) // 0644 is the permission to create if file does not exists
	CheckErrAndPanic(err)
}


func WriteDataAsJson(jsonPath string, jsonData []byte) {
	// sanity check
	fmt.Println(string(jsonData))
	// write to JSON file
	jsonFile, err := os.Create(jsonPath)
	CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	fmt.Println("JSON data written to ", jsonFile.Name())
}


func ReadFromFile(filePath string) []byte {
	// read the whole file at once
	data, err := ioutil.ReadFile(filePath)
	CheckErrAndPanic(err)
	return data
}


func CheckErrAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckErrToLog(err error) {
	CheckErrAndPanic(err) // TODO remove
	if err != nil {
		log.Fatal(err)
	}
}


func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}


func EncodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return string(pemEncodedPub)
}


func DecodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}


func DecodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}


func intMax(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}