package scripts

import (
	"encoding/json"
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
)

// create a basic node mapping, for simulating mixnet architecture.
func CreateNodesMap() {
	userAddresses := make(mixnet.UserAddressMapType, 0)
	userPublicKeys := make(mixnet.UserEncodedPublicKeyMapType, 0)
	userPrivateKeys := make(mixnet.UserEncodedPrivateKeyMapType, 0)
	for _, userName := range mixnet.UserNames {
		userAddresses[userName] = fmt.Sprintf(mixnet.AddressFormat, userName)
		privateKey, publicKey := mixnet.GenerateAsymmetricKeyPair()
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
