package main

import (
	"encoding/hex"
	"gonexa/account"
	"gonexa/nexa"
)

func main() {
	// NexaTest()
	// nexa.TxTest()
	// nexa.NexaTxHashTest()
	// nexa.NexaSignTxOneInputTest()
	// nexa.NexaSignTxTest()
	// nexa.TxID()
	pretxId := "1848efbaee543dd058d3529fcfe95652e149dcff02f820245426703c1c65a975"
	// pretxId = "3c6746ede34e13dce746c049ca7ec53ce63d95bb8b9cfa43e7d8f8e1ae4501de"
	txidReverse := nexa.HashReverseHex(pretxId)
	ouputIndexHex := nexa.TxAmountToLitteEndianHex(1)
	println("ouputIndexHex:", ouputIndexHex)
	preout := hex.EncodeToString(txidReverse) + ouputIndexHex
	// preout = pretxId + ouputIndexHex
	println("preout:", preout)
	// 16af7506d867a4f24da79ad4f70fc82bb4d3d9a79502e1f603d6096e25f0e8f2
	res := nexa.TxDoubleHash256(preout)
	println("res:", res)
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
