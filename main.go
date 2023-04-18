package main

import (
	"gonexa/account"
	"gonexa/nexa"
)

func main() {
	// NexaTest()
	// nexa.TxTest()
	// nexa.NexaTxHashTest()
	nexa.NexaSignTxOneInputTest()
	// nexa.NexaSignTxTest()
	// nexa.TxID()

}

func NexaTest() {

	account := account.GetMainAccount()
	// println(account.ToString())
	mnemonic := account.Mnenmonic
	// mnemonic = "champion turn truck wealth leave angry fatigue topic style core permit deny"
	pri, pub := nexa.NexaMnemonicToPrivateKeyAndPublicKey(mnemonic)
	println(pri, pub)
	nexa.NexaPublickeyToAddress(pub)
	pub, addr := nexa.NexaPrivateKeyToPublicKeyAndAddress(pri)
	println(pub, addr)
}
