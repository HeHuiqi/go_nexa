package nexa

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

func GetP() *big.Int {
	curveParams := btcec.Params()
	return curveParams.P
}
func GetN() *big.Int {
	curveParams := btcec.Params()
	return curveParams.N
}

func GetG() (*big.Int, *big.Int) {
	curveParams := btcec.Params()
	return curveParams.Gx, curveParams.Gy
}

func GetX() *big.Int {
	curveParams := btcec.Params()
	return curveParams.Gx
}
func GetY() *big.Int {
	curveParams := btcec.Params()
	return curveParams.Gx
}
