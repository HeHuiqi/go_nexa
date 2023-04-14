package main

import (
	"gonexa/account"
	"gonexa/nexa"
)

func main() {
	// NexaTest()

	nexa.TxTest()

}

func NexaTest() {

	account := account.GetMainAccount()
	// println(account.ToString())
	mnemonic := account.Mnenmonic
	pri, pub := nexa.NexaMnemonicToPrivateKeyAndPublicKey(mnemonic)
	println(pri, pub)

	nexa.NexaPublickeyToAddress(pub)
	pub, addr := nexa.NexaPrivateKeyToPublicKeyAndAddress(pri)
	println(pub, addr)
}
