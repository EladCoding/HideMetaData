package scripts

import (
	"github.com/EladCoding/HideMetaData/mixnet"
)


func RunNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}
