package nexa

import (
	"encoding/hex"

	"github.com/gcash/bchd/bchec"
)

func SchnorrSignMessageHash(msgHashHex string, pri string) (string, string) {
	privKeyBytes, _ := hex.DecodeString(pri)
	privKey, pubKey := bchec.PrivKeyFromBytes(bchec.S256(), privKeyBytes)
	hash, _ := hex.DecodeString(msgHashHex)
	signature, _ := privKey.SignSchnorr(hash)
	println("signature:", hex.EncodeToString(signature.Serialize()))
	return hex.EncodeToString(signature.Serialize()), hex.EncodeToString(pubKey.SerializeCompressed())

}
func ESDACSignMessageHash(msgHashHex string, pri string) (string, string) {
	privKeyBytes, _ := hex.DecodeString(pri)
	privKey, pubKey := bchec.PrivKeyFromBytes(bchec.S256(), privKeyBytes)
	hash, _ := hex.DecodeString(msgHashHex)
	signature, _ := privKey.SignECDSA(hash)
	println("signature:", hex.EncodeToString(signature.Serialize()))
	return hex.EncodeToString(signature.Serialize()), hex.EncodeToString(pubKey.SerializeCompressed())
}
