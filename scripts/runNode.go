package scripts

import (
	"github.com/EladCoding/HideMetaData/mixnet"
)

// Run one node at the mixnet architecture.
func RunNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}
