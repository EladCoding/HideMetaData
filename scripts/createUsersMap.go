package scripts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
)


func GenerateAsymmetricKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	CheckErrToLog(err)
	pubKey := privKey.PublicKey
	return privKey, pubKey
}


func CreateUsersMap() {
	userAddresses := make(UserAddressMap, 0)
	userPublicKeys := make(UserPublicKeyMap, 0)
	userPrivateKeys := make(UserPrivateKeyMap, 0)
	for _, userName := range UserNames {
		userAddresses[userName] = fmt.Sprintf(AddressFormat, userName)
		privateKey, publicKey := GenerateAsymmetricKeyPair()
		userPublicKeys[userName] = EncodePublicKey(&publicKey)
		userPrivateKeys[userName] = EncodePrivateKey(privateKey)
	}
	userAddressesjsonData, err := json.MarshalIndent(userAddresses, "", "\t")
	CheckErrAndPanic(err)
	WriteToFile(UserAddressesMapPath, userAddressesjsonData)
	userPublicKeysjsonData, err := json.MarshalIndent(userPublicKeys, "", "\t")
	CheckErrAndPanic(err)
	WriteToFile(UserPublicKeysMapPath, userPublicKeysjsonData)
	userPrivateKeysjsonData, err := json.MarshalIndent(userPrivateKeys, "", "\t")
	CheckErrAndPanic(err)
	WriteToFile(UserPrivateKeysMapPath, userPrivateKeysjsonData)
}
