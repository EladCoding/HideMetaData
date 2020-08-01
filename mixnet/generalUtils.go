package mixnet

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

type UserAddressMapType map[string]string
type UserEncodedPublicKeyMapType map[string]string
type UserDecodedPublicKeyMapType map[string]*ecdsa.PublicKey
type UserEncodedPrivateKeyMapType map[string]string
type UserDecodedPrivateKeyMapType map[string]*ecdsa.PrivateKey
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

// Write data to a given file.
func WriteToFile(filePath string, data []byte) {
	// write the whole body at once
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	err := ioutil.WriteFile(filePath, data, 0644) // 0644 is the permission to create if file does not exists
	CheckErrAndPanic(err)
}

// Write data to a given file, as a json file.
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

// Read data to a given file.
func ReadFromFile(filePath string) []byte {
	// read the whole file at once
	data, err := ioutil.ReadFile(filePath)
	CheckErrAndPanic(err)
	return data
}

// Check an error message, and panic if an error discovered.
func CheckErrAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Check an error message, and write the error to log if discovered.
func CheckErrToLog(err error) {
	if err != nil {
		log.Fatal(err)
	}
	CheckErrAndPanic(err)
}

// Check if a slice contain a given string.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


func intMax(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}


func GetMemUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / float64(m.Sys)

}


func InitLogFile() {
	t := time.Now()
	logFileDir := "./logs/"
	ValidateDirExists(logFileDir)
	logFileName := logFileDir + t.Format("20060102150405") + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	log.Printf("Initialize log file")
}


func ValidateDirExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0700)
	}
}

