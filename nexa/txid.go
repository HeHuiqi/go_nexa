package nexa

func NexaTxIdemHex(inputs []NexaInputOutpoint, outputs []NexaOutput, lockTime uint32) string {
	txIdemHex := TxVersion()
	inputCount := uint64(len(inputs))
	txIdemHex += VarIntToHex(inputCount)
	for i := uint64(0); i < inputCount; i++ {
		input := inputs[i]
		txIdemHex += input.ToIdemHexString()
	}

	outputCount := uint64(len(outputs))
	txIdemHex += VarIntToHex(outputCount)
	for i := uint64(0); i < outputCount; i++ {
		output := outputs[i]
		txIdemHex += output.ToHexString()
	}
	txIdemHex += TxLocktime(lockTime)
	return txIdemHex
}
func NexaTxToSatisfierHex(inputs []NexaInputOutpoint) string {
	inputCount := uint32(len(inputs))
	satisfierHex := Int32ToLitteEndianHex(inputCount)
	for i := uint32(0); i < inputCount; i++ {
		input := inputs[i]
		//输入的签名脚本
		satisfierHex += input.SignatureScript
		// Opcode.map.OP_INVALIDOPCODE) = 0xff
		satisfierHex += "ff"
	}
	return satisfierHex
}

func NexaTxIdemHash(idemHex string) string {
	return TxDoubleHash256(idemHex)
}
func NexaTxSatisfierHash(satisfierHex string) string {
	return TxDoubleHash256(satisfierHex)
}
func NexaTxIdHash(idemHex string, satisfierHex string) string {
	idHex := NexaTxIdemHash(idemHex) + NexaTxSatisfierHash(satisfierHex)
	hashHex := TxDoubleHash256(idHex)
	return hashHex
}

// 注意此方法需要在签名交易后调用
func NexaTxIdAndTxIdem(inputs []NexaInputOutpoint, outputs []NexaOutput, lockTime uint32) (txId string, txIdem string) {
	txIdemHexHex := NexaTxIdemHex(inputs, outputs, lockTime)
	// println("txIdemHexHex:", txIdemHexHex)
	txSatisfierHex := NexaTxToSatisfierHex(inputs)
	// println("txSatisfierHex:", txSatisfierHex)
	txIdemHash := NexaTxIdemHash(txIdemHexHex)
	txIdHash := NexaTxIdHash(txIdemHexHex, txSatisfierHex)
	txId = HashReverseHex(txIdHash)
	txIdem = HashReverseHex(txIdemHash)
	// println("txId:", txId)
	// println("txIdem:", txIdem)
	return
}

/*


Transaction Idem Calculation
Serialize the following transaction fields using standard bitcoin serialization algorithms:

1 version
2 inputs
	prevout
	sequence
	amount
	NOTE: the satisfier script (scriptSig) is not serialized – not even as an empty array
3 outputs
4 locktime

Transaction Id Calculation
1 Create the “satisfiersHash” by double SHA256 hashing the following byte stream:
	number of inputs as a little endian 4 byte number
	for each input:
	satisfier script (script sig)
	0xFF
2 Calculate the transaction Idem
3 Concatenate the Idem with the satisfiersHash.
4 The transaction Id is the double SHA256 of the result of step 3.
*/
