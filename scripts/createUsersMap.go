package scripts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
)


func GenerateAsymmetricKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	mixnet.CheckErrToLog(err)
	pubKey := privKey.PublicKey
	return privKey, pubKey
}


func CreateUsersMap() {
	userAddresses := make(mixnet.UserAddressMap, 0)
	userPublicKeys := make(mixnet.UserPublicKeyMap, 0)
	userPrivateKeys := make(mixnet.UserPrivateKeyMap, 0)
	for _, userName := range mixnet.UserNames {
		userAddresses[userName] = fmt.Sprintf(mixnet.AddressFormat, userName)
		privateKey, publicKey := GenerateAsymmetricKeyPair()
		userPublicKeys[userName] = mixnet.EncodePublicKey(&publicKey)
		userPrivateKeys[userName] = mixnet.EncodePrivateKey(privateKey)
	}
	userAddressesjsonData, err := json.MarshalIndent(userAddresses, "", "\t")
	mixnet.CheckErrAndPanic(err)
	mixnet.WriteToFile(mixnet.UserAddressesMapPath, userAddressesjsonData)
	userPublicKeysjsonData, err := json.MarshalIndent(userPublicKeys, "", "\t")
	mixnet.CheckErrAndPanic(err)
	mixnet.WriteToFile(mixnet.UserPublicKeysMapPath, userPublicKeysjsonData)
	userPrivateKeysjsonData, err := json.MarshalIndent(userPrivateKeys, "", "\t")
	mixnet.CheckErrAndPanic(err)
	mixnet.WriteToFile(mixnet.UserPrivateKeysMapPath, userPrivateKeysjsonData)
}
