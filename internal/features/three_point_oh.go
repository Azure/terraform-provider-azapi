package features

import (
	"os"
	"strings"
)

func ThreePointOh() bool {
	return strings.EqualFold(os.Getenv("ARM_THREEPOINTZERO_BETA"), "true")
}
