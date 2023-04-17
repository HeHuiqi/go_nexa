package nexa

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"gonexa/account"
)

// 已废弃
func TxPrevOutHash(outPointHex string, inputIndex uint8) string {

	outPoint := TxOutPoint(outPointHex, inputIndex)
	hasResult := TxDoubleHash256(outPoint)
	println("PrevOutHash:", hasResult)
	return hasResult

}
func TxSequenceHash(sequence uint32) string {

	wbuf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(wbuf, binary.LittleEndian, sequence)
	sequenceHex := hex.EncodeToString(wbuf.Bytes())
	hasResult := TxDoubleHash256(sequenceHex)
	println("TxTxSequenceHash:", hasResult)
	return hasResult

}

func TxAmountHash(amount uint64) string {
	amoutHex := TxAmountToLitteEndianHex(amount)
	hasResult := TxDoubleHash256(amoutHex)
	println("TxAmountHash:", hasResult)
	return hasResult

}

func NexaOuputsSerialize(sendAmount uint64, toAddress string, changeAmount uint64, changeAddress string) string {
	sendHex := NexaOuputSerialize(1, sendAmount, toAddress)
	changeHex := NexaOuputSerialize(1, changeAmount, changeAddress)
	ret := sendHex + changeHex
	println("NexaOuputsSerialize:", ret)
	return ret
}
func NexaOuputsHash(ouputsHex string) string {
	ret := TxDoubleHash256(ouputsHex)
	println("NexaOuputsHash:", ret)
	return ret
}

func NexaScriptSerialize(signType uint8, scriptHex string) string {

	// signType = 0 表示 NexaSignTypeAll
	if signType == 0 {
		// 02 len
		//  Opcode.OP_FROMALTSTACK= 108 = 6c, Opcode.OP_CHECKSIGVERIFY= 173 = 0xad
		return "026cad"
	}
	lenHex := Int8ToHexString(int8(len(scriptHex)/2 + 1))
	return lenHex + scriptHex + "6cad"
}
func NexaP2PKTScriptSerializeSignTypeAll() string {
	ret := NexaScriptSerialize(0, "")
	println("NexaP2PKTScriptSerializeSignTypeAll:", ret)
	return ret
}

func NeaxSignTypeAllHex() string {
	return "00"
}

func NexaTx(outpointHex string, outputIdx uint8, ouputBalance uint64, sequence uint32,
	fromAddr string,
	changeAddr string,
	changeAmount uint64,
	toAddr string,
	sendAmount uint64,
	locktime uint32,
	privateKey string,
	pubkeyHex string,
) {
	hashBuf := ""
	version := TxVersion()
	hashBuf += version
	hashBuf += TxPrevOutHash(outpointHex, outputIdx)
	hashBuf += TxAmountHash(ouputBalance)
	hashBuf += TxSequenceHash(sequence)
	hashBuf += NexaP2PKTScriptSerializeSignTypeAll()

	outputs := NexaOuputsSerialize(sendAmount, toAddr, changeAmount, changeAddr)
	hashBuf += NexaOuputsHash(outputs)

	locktimeHex := TxLocktime(locktime)
	hashBuf += locktimeHex
	sinTypeHex := NeaxSignTypeAllHex()
	hashBuf += sinTypeHex

	println("hashBuf:", hashBuf)

	hashHex := TxDoubleHash256(hashBuf)
	println("hashHex=", hashHex)

	// hashHexBytes, _ := hex.DecodeString(hashHex)
	// reverseHashBytes := HashReverse(hashHexBytes)
	// reverseHashHex := hex.EncodeToString(reverseHashBytes)
	// println("reverseHashHex=", reverseHashHex)

	sign, _ := NexaSign(hashHex, privateKey)
	println("sign:", sign)

	/*
		https://spec.nexa.org/protocol/blockchain/transaction 交易结构

		00 version
		01 input_count
		00 output_index
		94a577802f958a3c151411f91ce434897b705964ed5ae56bbf26df0d314bc5b3 reverse(outpoint)
		64 sing_all_len
		2221 pubkey_all_len
		02c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71 pubkey
		40 sign_len
		22d5f6fdd0c63d51f2f9ac9c3e02f8eb7cb1012637cd89788ec2fc62cb325861
		a127ed99995d0272ae92a793659e1cb0515031f0c2c70dc47ed5a061e4bf9dd8
		feffffff sequence
		1027000000000000 input_amount
		02 output_count
		01 output_type
		a00f000000000000  output_amount1
		170051141476bc2644608f984125fd4b7875f067bc5953b8 output_script_pubkey1
		01 output_type
		df14000000000000 output_amount2
		170051141129a5ec6501c423c686247dfe7f413b4ebf7449 output_script_pubkey2
		f6dc0300 locktime = blockHeight


		00
		01
		00
		94a577802f958a3c151411f91ce434897b705964ed5ae56bbf26df0d314bc5b3
		64
		2221
		02c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71
		40
		22d5f6fdd0c63d51f2f9ac9c3e02f8eb7cb1012637cd89788ec2fc62cb325861
		a127ed99995d0272ae92a793659e1cb0515031f0c2c70dc47ed5a061e4bf9dd8
		feffffff
		1027000000000000
		02
		01
		a00f000000000000
		170051141476bc2644608f984125fd4b7875f067bc5953b8
		01
		df14000000000000
		170051141129a5ec6501c423c686247dfe7f413b4ebf7449
		f6dc0300

	*/

	txRawStr := version
	input_count := Int8ToLittleEndianHexString(1)
	txRawStr += input_count

	txOutPointHex := TxOutPoint(outpointHex, outputIdx)
	txRawStr += txOutPointHex
	pubkeyFormat := FormatData(pubkeyHex)
	pubkeyFormat = FormatData(pubkeyFormat)
	// println("pubkeyFormat:", pubkeyFormat)
	signFromat := FormatData(sign)
	// println("signFromat:", signFromat)

	signAllFormat := FormatData(pubkeyFormat + signFromat)
	// println("signAllFormat:", signAllFormat)
	txRawStr += signAllFormat

	sequenceHex := Int32ToLitteEndianHex(sequence)
	txRawStr += sequenceHex
	inputAmoutHex := TxAmountToLitteEndianHex(ouputBalance)
	txRawStr += inputAmoutHex
	// outputs
	output_count := Int8ToLittleEndianHexString(2)
	txRawStr += output_count

	output_type := "01"
	txRawStr += output_type
	sendAmountHex := TxAmountToLitteEndianHex(sendAmount)
	txRawStr += sendAmountHex
	toScript, _ := NexaP2KTScript(toAddr)
	txRawStr += toScript

	txRawStr += output_type
	changeAmountHex := TxAmountToLitteEndianHex(changeAmount)
	txRawStr += changeAmountHex
	changeScript, _ := NexaP2KTScript(changeAddr)
	txRawStr += changeScript

	txRawStr += locktimeHex

	println("txRawStr:", txRawStr)

}

func TxTest() {
	// https://explorer.nexa.org/tx/a744499844b2276b41a667d0efcead5ebfceb512b3c7dac7c4ca184acb552d5d
	// outpointHex 是返回的utxo中的 outpoint 字段的值，不是txid
	outpointHex := "b3c54b310ddf26bf6be55aed6459707b8934e41cf91114153c8a952f8077a594"
	outputIdx := uint8(0)
	ouputBalance := uint64(10000)
	sequence := 0xfffffffe
	fromAccount := account.GetMainAccount()
	changeAmount := uint64(0x14df)
	// toAccount := account.GetAccount(1)
	// toAddress := toAccount.Address
	toAddress := "nexa:nqtsq5g5z3mtcfjyvz8essf9l49hsa0sv779j5acw6sdj4e8"
	sendAmount := uint64(0x0fa0)
	locktime := uint32(253174) // future block height

	NexaTx(outpointHex, outputIdx, ouputBalance, uint32(sequence),
		fromAccount.Address, fromAccount.ChangeAddress, changeAmount,
		toAddress, sendAmount, locktime, fromAccount.PrivateKey, fromAccount.PublicKey)

}
