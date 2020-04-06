package mixnet

import (
	"fmt"
	"log"
	"os"
)

//// types
type userInfoMap map[string][2]string

//// vars
var working_dir, _ = os.Getwd()
// cipher vars
var RsaKeyBits = 2048
var CipherRsaLen = RsaKeyBits / 8
var AesKeyBytes = 32

var PathLen = 3
// servers map
var ServerPublicKeyPathFormat = working_dir + "/keys/server%s/public_key.txt"
var ServerPrivateKeyPathFormat = working_dir + "/keys/server%s/private_key.txt"
var ServerHost = "localhost"
var ServerPortFormat = "9%s"
var ServerAddressFormat = ServerHost + ":" + ServerPortFormat
// mediators map
var MediatorPublicKeyPathFormat = working_dir + "/keys/mediator%s/public_key.txt"
var MediatorPrivateKeyPathFormat = working_dir + "/keys/mediator%s/private_key.txt"
var MediatorHost = "localhost"
var MediatorHostFormat = "8%s"
var MediatorAddressFormat = MediatorHost + ":" + MediatorHostFormat
// clients map
var ClientPublicKeyPathFormat = working_dir + "/keys/client%s/public_key.txt"
var ClientPrivateKeyPathFormat = working_dir + "/keys/client%s/private_key.txt"
var ClientHost = "localhost"
var ClientHostFormat = "7%s"
var ClientAddressFormat = ClientHost + ":" + ClientHostFormat

var PublicKeyPathSpot = 0
var AddressSpot = 1
var UserNameLen = 3

var usersMap = userInfoMap{
	"001": {fmt.Sprintf(ServerPublicKeyPathFormat, "001"), fmt.Sprintf(ServerAddressFormat, "001")},
	"002": {fmt.Sprintf(ServerPublicKeyPathFormat, "002"), fmt.Sprintf(ServerAddressFormat, "002")},
	"003": {fmt.Sprintf(ServerPublicKeyPathFormat, "003"), fmt.Sprintf(ServerAddressFormat, "003")},
	"101": {fmt.Sprintf(MediatorPublicKeyPathFormat, "101"), fmt.Sprintf(MediatorAddressFormat, "101")},
	"102": {fmt.Sprintf(MediatorPublicKeyPathFormat, "102"), fmt.Sprintf(MediatorAddressFormat, "102")},
	"103": {fmt.Sprintf(MediatorPublicKeyPathFormat, "103"), fmt.Sprintf(MediatorAddressFormat, "103")},
	"201": {fmt.Sprintf(ClientPublicKeyPathFormat, "201"), fmt.Sprintf(ClientAddressFormat, "201")},
	"202": {fmt.Sprintf(ClientPublicKeyPathFormat, "202"), fmt.Sprintf(ClientAddressFormat, "202")},
	"203": {fmt.Sprintf(ClientPublicKeyPathFormat, "203"), fmt.Sprintf(ClientAddressFormat, "203")},
}

//// funcs
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
