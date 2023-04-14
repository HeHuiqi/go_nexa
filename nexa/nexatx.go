package nexa

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"gonexa/account"

	"github.com/gcash/bchd/chaincfg/chainhash"
)

func Int8ToHexString(intNum int8) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 2))
	binary.Write(wbuf, binary.BigEndian, intNum)
	return hex.EncodeToString(wbuf.Bytes())
}

func Int8ToLittleEndianHexString(intNum int8) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 2))
	binary.Write(wbuf, binary.LittleEndian, intNum)
	return hex.EncodeToString(wbuf.Bytes())
}

func TxVersion() string {
	return "00"
}

func TxOutPoint(preOutPointHex string, outIndex uint8) string {
	utxoTxHash, _ := chainhash.NewHashFromStr(preOutPointHex)
	wbuf := bytes.NewBuffer(make([]byte, 0, len(utxoTxHash)+2))
	binary.Write(wbuf, binary.LittleEndian, outIndex)
	binary.Write(wbuf, binary.LittleEndian, utxoTxHash)
	outPoint := hex.EncodeToString(wbuf.Bytes())
	println("outPoint:", outPoint)
	return outPoint
}

func TxAmountToLitteEndianHex(amount uint64) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(wbuf, binary.LittleEndian, amount)
	amountStr := hex.EncodeToString(wbuf.Bytes())
	// println("amountStr:", amountStr)
	return amountStr
}
func Int32ToLitteEndianHex(num uint32) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 4))
	binary.Write(wbuf, binary.LittleEndian, num)
	numStr := hex.EncodeToString(wbuf.Bytes())
	// println("Int32ToLitteEndianHex:", numStr)
	return numStr
}

func TxSequence() string {
	// 0xfffffffe little end
	return "feffffff"
}

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

func TxDoubleHash256(hexStr string) string {
	sequenceBytes, _ := hex.DecodeString(hexStr)
	doubleHash := chainhash.DoubleHashB(sequenceBytes)
	hasResult := hex.EncodeToString(doubleHash)
	return hasResult
}

func NexaP2KTScript(address string) (string, string) {
	result, prefix, _, _ := CheckDecodeCashAddress(address)
	ret := hex.EncodeToString(result)
	// println(ret, prefix)
	return ret, prefix
}

func NexaOuputSerialize(outputType uint8, outputAmount uint64, toAddress string) string {
	typeHex := Int8ToHexString(int8(outputType))
	outputScpritHex, _ := NexaP2KTScript(toAddress)
	amountHex := TxAmountToLitteEndianHex(outputAmount)
	ret := typeHex + amountHex + outputScpritHex
	// println("NexaOuputSerialize:", ret)
	return ret
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

func TxLocktime(blockHeght uint32) string {
	ret := Int32ToLitteEndianHex(blockHeght)
	println("TxLocktime:", ret)
	return ret
}
func NeaxSinTypeAllHex() string {
	return "00"
}

func NexaSign(msgHash string, priHex string) string {
	priBytes, _ := hex.DecodeString(priHex)
	msgBytes, _ := hex.DecodeString(msgHash)
	ret := Signature(priBytes, msgBytes)
	return ret
}

func FormatData(hexStr string) string {
	hexStrFormat := Int8ToHexString(int8(len(hexStr)/2)) + hexStr
	// println("hexStrFormat:", hexStrFormat)
	return hexStrFormat
}

func TxTest() {
	// https://explorer.nexa.org/tx/a744499844b2276b41a667d0efcead5ebfceb512b3c7dac7c4ca184acb552d5d
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
	locktime := uint32(253174)

	NexaTx(outpointHex, outputIdx, ouputBalance, uint32(sequence),
		fromAccount.Address, fromAccount.ChangeAddress, changeAmount,
		toAddress, sendAmount, locktime, fromAccount.PrivateKey, fromAccount.PublicKey)

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
	sinTypeHex := NeaxSinTypeAllHex()
	hashBuf += sinTypeHex

	println("hashBuf:", hashBuf)

	hashHex := TxDoubleHash256(hashBuf)
	println("hashHex=", hashHex)

	// hashHexBytes, _ := hex.DecodeString(hashHex)
	// reverseHashBytes := HashReverse(hashHexBytes)
	// reverseHashHex := hex.EncodeToString(reverseHashBytes)
	// println("reverseHashHex=", reverseHashHex)

	sign := NexaSign(hashHex, privateKey)
	println("sign:", sign)
	// cde47ad66247351daa92597537989bf6902183c8cc974825f8023077f00637ef350df588341a110e693d3a766b796b74d09b3d0a8543ed1b538f94297b972b82

	/*

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
		f6dc03 locktime = blockHeight
		00 sigtype


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
		f6dc03
		00

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
