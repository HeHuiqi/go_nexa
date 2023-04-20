package nexa

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"gonexa/account"

	"github.com/gcash/bchd/chaincfg/chainhash"
)

func HashReverseHex(hexStr string) string {
	inBytes, _ := hex.DecodeString(hexStr)
	return hex.EncodeToString(HashReverse(inBytes))
}
func HashReverse(inBytes []byte) []byte {
	HashSize := len(inBytes)
	dst := make([]byte, HashSize)
	for i, b := range inBytes[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = inBytes[HashSize-1-i], b
	}
	return dst[:]
}

func TxVersion() string {
	return "00"
}
func TxSequence(sequence uint32) string {

	wbuf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(wbuf, binary.LittleEndian, sequence)
	sequenceHex := hex.EncodeToString(wbuf.Bytes())
	return sequenceHex
	// 0xfffffffe little end
	// return "feffffff"
}

func Int8ToHexString(intNum int8) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 2))
	binary.Write(wbuf, binary.BigEndian, intNum)
	return hex.EncodeToString(wbuf.Bytes())
}

func Int8ToLittleEndianHexString(intNum uint8) string {
	wbuf := bytes.NewBuffer(make([]byte, 0, 2))
	binary.Write(wbuf, binary.LittleEndian, intNum)
	return hex.EncodeToString(wbuf.Bytes())
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

func FormatData(hexStr string) string {
	hexStrFormat := Int8ToHexString(int8(len(hexStr)/2)) + hexStr
	// println("hexStrFormat:", hexStrFormat)
	return hexStrFormat
}
func TxDoubleHash256(hexStr string) string {
	sequenceBytes, _ := hex.DecodeString(hexStr)
	doubleHash := chainhash.DoubleHashB(sequenceBytes)
	hasResult := hex.EncodeToString(doubleHash)
	return hasResult
}

// index + preOutPointHex
func TxOutPoint(preOutPointHex string, outIndex uint8) string {
	utxoTxHash, _ := chainhash.NewHashFromStr(preOutPointHex)
	wbuf := bytes.NewBuffer(make([]byte, 0, len(utxoTxHash)+2))
	binary.Write(wbuf, binary.LittleEndian, outIndex)
	binary.Write(wbuf, binary.LittleEndian, utxoTxHash)
	outPoint := hex.EncodeToString(wbuf.Bytes())
	println("outPoint:", outPoint)
	return outPoint
}

type NexaInputOutpoint struct {
	// outpointHex 是返回的utxo中的 outpoint 字段的值，不是txid
	OutpointHex     string
	InputType       uint8 //必须是0
	OutputAmount    uint64
	Sequence        uint32
	SignatureScript string
}

func (input *NexaInputOutpoint) ToHexString() string {
	return TxOutPoint(input.OutpointHex, input.InputType)
}
func (input *NexaInputOutpoint) ToIdemHexString() string {
	idemHex := Int8ToHexString(int8(input.InputType))
	idemHex += HashReverseHex(input.OutpointHex)
	idemHex += TxSequence(input.Sequence)
	idemHex += TxAmountToLitteEndianHex(input.OutputAmount)
	return idemHex
}
func (input *NexaInputOutpoint) SetSignatureScript(sigScript string) {
	input.SignatureScript = sigScript
}

func NewInputOutpoint(outpointHex string, inputType uint8, ouputBalance uint64, sequence uint32) NexaInputOutpoint {

	return NexaInputOutpoint{OutpointHex: outpointHex, InputType: inputType, OutputAmount: ouputBalance, Sequence: sequence}
}

func NexaP2KTScript(address string) (string, string) {
	result, prefix, _, _ := CheckDecodeCashAddress(address)
	ret := hex.EncodeToString(result)
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

func NexaInputsToHash(inputs []NexaInputOutpoint) string {
	hex := inputs[0].ToHexString()
	for i := 1; i < len(inputs); i++ {
		hex += inputs[i].ToHexString()
	}
	hash := TxDoubleHash256(hex)
	return hash
}
func NexaInputsSequenceToHash(inputs []NexaInputOutpoint) string {
	hex := TxSequence(inputs[0].Sequence)
	for i := 1; i < len(inputs); i++ {
		hex += TxSequence(inputs[i].Sequence)
	}
	hash := TxDoubleHash256(hex)
	return hash
}

func NexaInputsAmountToHash(inputs []NexaInputOutpoint) string {
	hex := TxAmountToLitteEndianHex(inputs[0].OutputAmount)
	for i := 1; i < len(inputs); i++ {
		hex += TxAmountToLitteEndianHex(inputs[i].OutputAmount)
	}
	hash := TxDoubleHash256(hex)
	return hash
}

type NexaOutput struct {
	OutputType   uint8
	OutputAmount uint64
	OutputScript string
	Address      string
}

func NexaNewOutput(outputType uint8, amount uint64, address string) (*NexaOutput, error) {
	script, _, _, err := CheckDecodeCashAddress(address)
	if err != nil {
		return nil, err
	}
	return &NexaOutput{OutputType: outputType, OutputAmount: amount, Address: address, OutputScript: hex.EncodeToString(script)}, nil
}
func (output *NexaOutput) ToHexString() string {
	hex := Int8ToLittleEndianHexString(output.OutputType)
	hex += TxAmountToLitteEndianHex(output.OutputAmount)
	hex += output.OutputScript
	return hex
}

func NexaOutpusToHash(outputs []NexaOutput) string {

	hex := outputs[0].ToHexString()
	for i := 1; i < len(outputs); i++ {
		hex += outputs[i].ToHexString()
	}
	hash := TxDoubleHash256(hex)
	return hash
}

func NexaSignTypeHex(signType uint8) string {
	if signType == 0 {
		return "00"
	}
	return Int8ToLittleEndianHexString(signType)
}
func TxLocktime(blockHeght uint32) string {
	ret := Int32ToLitteEndianHex(blockHeght)
	println("TxLocktime:", ret)
	return ret
}

func NexaTxHash(inputs []NexaInputOutpoint, outputs []NexaOutput, lockTime uint32, signType uint8) string {
	hashBuf := ""
	version := TxVersion()
	hashBuf += version
	inputsHash := NexaInputsToHash(inputs)
	println("inputsHash:", inputsHash)
	hashBuf += inputsHash

	inputsAmountHash := NexaInputsAmountToHash(inputs)
	println("inputsAmountHash:", inputsAmountHash)
	hashBuf += inputsAmountHash

	inputsSequenceHash := NexaInputsSequenceToHash(inputs)
	println("inputsSequenceHash:", inputsSequenceHash)
	hashBuf += inputsSequenceHash

	scriptHex := NexaP2PKTScriptSerializeSignTypeAll()
	println("scriptHex:", scriptHex)
	hashBuf += scriptHex

	outpusHash := NexaOutpusToHash(outputs)
	println("outpusHash:", outpusHash)
	hashBuf += outpusHash

	lockTimeHex := TxLocktime(lockTime)
	hashBuf += lockTimeHex
	sigtypeHex := NexaSignTypeHex(signType)
	hashBuf += sigtypeHex

	println("hashBuf:", hashBuf)
	msgHash := TxDoubleHash256(hashBuf)
	println("msgHash:", msgHash)
	return msgHash
}
func NexaSign(msgHash string, priHex string) (string, string) {
	// priBytes, _ := hex.DecodeString(priHex)
	// msgBytes, _ := hex.DecodeString(msgHash)
	// ret, pubHex := Signature(priBytes, msgBytes)
	ret, pubHex := SchnorrSignMessageHash(msgHash, priHex)
	return ret, pubHex
}
func FormatSignRaw(signHex string, pubHex string) string {
	pubFormat := FormatData(pubHex)
	pubFormat = FormatData(pubFormat)
	signFormat := FormatData(signHex)

	ret := FormatData(pubFormat + signFormat)
	return ret
}
func FormatSignRawToSigScript(signHex string, pubHex string) string {
	pubFormat := FormatData(pubHex)
	pubFormat = FormatData(pubFormat)
	signFormat := FormatData(signHex)
	return pubFormat + signFormat
}
func NexaSignTx(inputs []NexaInputOutpoint, outputs []NexaOutput, lockTime uint32, msgHashHex string, priHex string) string {
	signHex, pubHex := NexaSign(msgHashHex, priHex)
	println("signHex:", signHex)
	println("pubHex:", pubHex)
	signFormat := FormatSignRaw(signHex, pubHex)

	inputsFormat := TxVersion()
	inputCount := uint8(len(inputs))
	inputCountHex := Int8ToLittleEndianHexString(inputCount)
	inputsFormat += inputCountHex

	for i := uint8(0); i < inputCount; i++ {
		input := inputs[i]
		inputsFormat += input.ToHexString()
		inputsFormat += signFormat
		inputsFormat += TxSequence(input.Sequence)
		inputsFormat += TxAmountToLitteEndianHex(input.OutputAmount)

	}
	println("inputsFormat:", inputsFormat)
	outpuCount := uint8(len(outputs))
	outpuCountHex := Int8ToLittleEndianHexString(outpuCount)
	outputsFormat := outpuCountHex
	for i := uint8(0); i < outpuCount; i++ {
		output := outputs[i]
		outputsFormat += NexaOuputSerialize(output.OutputType, output.OutputAmount, output.Address)
	}
	println("outputsFormat:", outputsFormat)

	ret := inputsFormat + outputsFormat
	ret += TxLocktime(lockTime)

	//设置输入的解锁脚本
	for i := 0; i < len(inputs); i++ {
		inputs[i].SignatureScript = FormatSignRawToSigScript(signHex, pubHex)
	}

	return ret
}

func NexaTxHashTest() {
	//
	inputType := uint8(0)
	input1 := NewInputOutpoint("b3c54b310ddf26bf6be55aed6459707b8934e41cf91114153c8a952f8077a594", inputType, 10000, 0xfffffffe)
	input2 := NewInputOutpoint("84f6adb5ad2b1af7ff3026d16843cc123fc260f2ab4c3cd75b2d20df1dc431e4", inputType, 13000, 0xfffffffe)
	inputs := []NexaInputOutpoint{input1, input2}

	output1, _ := NexaNewOutput(1, uint64(0x4e20), "nexa:nqtsq5g5z3mtcfjyvz8essf9l49hsa0sv779j5acw6sdj4e8")
	output2, _ := NexaNewOutput(1, uint64(0x0771), "nexa:nqtsq5g5zy56tmr9q8zz835xy37lul6p8d8t7azfpuz2gs4e")
	outputs := []NexaOutput{*output1, *output2}
	lockTime := uint32(253841) //目前使用未来的一个区块高度
	//signType = 0 表示 NexaSignTypeAll 一般都用这个
	signType := uint8(0)

	msgHash := NexaTxHash(inputs, outputs, lockTime, signType)
	println("msgHash:", msgHash)

}

func NexaSignTxOneInputTest() {
	// https://explorer.nexa.org/tx/1848efbaee543dd058d3529fcfe95652e149dcff02f820245426703c1c65a975
	// outpointHex 是返回的utxo中的 outpoint 字段的值，不是txid
	inputType := uint8(0) //必须是0

	inputAmount := uint64(130 * 100)
	feeAmount := uint64(8 * 100)
	sendAmount := uint64(10 * 100)
	changeAmount := inputAmount - sendAmount - feeAmount
	outpointHex := "b7ec98cac00ad2533175632be7885b3433b1fa6a3a2b4076e23c60e344a6ec34"
	input1 := NewInputOutpoint(outpointHex, inputType, inputAmount, 0xfffffffe)
	inputs := []NexaInputOutpoint{input1}

	output1, _ := NexaNewOutput(1, sendAmount, "nexa:nqtsq5g5z3mtcfjyvz8essf9l49hsa0sv779j5acw6sdj4e8")
	output2, _ := NexaNewOutput(1, changeAmount, "nexa:nqtsq5g50j2ggw07vhj3rdmn0lnfh6fp6ef3gy47qudau8mx")
	outputs := []NexaOutput{*output1, *output2}

	//目前使用未来的一个区块高度
	// lockTime := uint32(253841)
	lockTime := uint32(255888)
	// lockTime := uint32(255899)
	signType := uint8(0)
	msgHash := NexaTxHash(inputs, outputs, lockTime, signType)

	priHex := account.GetMainAccount().PrivateKey
	signTxRaw := NexaSignTx(inputs, outputs, lockTime, msgHash, priHex)
	println("signTxRaw:", signTxRaw)

	txId, txIdem := NexaTxIdAndTxIdem(inputs, outputs, lockTime)
	println("txId:", txId)
	println("txIdem:", txIdem)

}

func NexaSignTxTest() {
	// outpointHex 是返回的utxo中的 outpoint 字段的值，不是txid
	inputType := uint8(0) //必须是0

	input1 := NewInputOutpoint("b3c54b310ddf26bf6be55aed6459707b8934e41cf91114153c8a952f8077a594", inputType, 10000, 0xfffffffe)
	input2 := NewInputOutpoint("84f6adb5ad2b1af7ff3026d16843cc123fc260f2ab4c3cd75b2d20df1dc431e4", inputType, 13000, 0xfffffffe)
	inputs := []NexaInputOutpoint{input1, input2}

	output1, _ := NexaNewOutput(1, uint64(9000), "nexa:nqtsq5g5z3mtcfjyvz8essf9l49hsa0sv779j5acw6sdj4e8")
	output2, _ := NexaNewOutput(1, uint64(13000), "nexa:nqtsq5g50j2ggw07vhj3rdmn0lnfh6fp6ef3gy47qudau8mx")
	outputs := []NexaOutput{*output1, *output2}

	//目前使用未来的一个区块高度
	// lockTime := uint32(253841)
	lockTime := uint32(255219)
	signType := uint8(0)
	msgHash := NexaTxHash(inputs, outputs, lockTime, signType)

	priHex := account.GetMainAccount().PrivateKey
	signTxRaw := NexaSignTx(inputs, outputs, lockTime, msgHash, priHex)
	println("signTxRaw:", signTxRaw)

	txId, txIdem := NexaTxIdAndTxIdem(inputs, outputs, lockTime)
	println("txId:", txId)
	println("txIdem:", txIdem)

}
