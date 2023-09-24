package share

import (
	"fmt"
	"math/big"
	"testing"
)

func TestName(t *testing.T) {

	secret := big.NewInt(2)

	shares := ShareSecret(secret, 3, 4)
	fmt.Println("shares: ", shares)

	secret1 := RecoverSecret(shares.ShareValues...)
	fmt.Println("recovered secret: ", secret1)
}
