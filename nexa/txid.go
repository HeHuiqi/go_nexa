package nexa

/*

"txid": "4a5733d194cd9572937b5ef766c35c631301430251f42d24ab343ec150478481",
"txidem": "a87876c510a3823c041db9a04c6925014b8bd82a91862e721b4149d70d5a25c5",


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


22
21
02c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71
40
3c9656f19082e40e5a212efeeb29a640f05dc8bd7719e3d817272f659eed3ce2
9fb2492feee60bde5720c5c9cdd62e62fa295c3eef8dba618f0e01e408930ea3
*/

func BtcTxSerialize(inputs []NexaInputOutpoint, outputs []NexaOutput, lockTime uint32, signType uint8) {
	version := TxVersion()
	ret := ""
	ret += version

	println("BtcTxSerialize:", ret)

}

func TxID() {
	println("-----TxID-----")
	/*
		"txid": "4a5733d194cd9572937b5ef766c35c631301430251f42d24ab343ec150478481",
		"txidem": "a87876c510a3823c041db9a04c6925014b8bd82a91862e721b4149d70d5a25c5",

	*/
	inputCountHex := "02"
	inputsScriptHex := "64222102c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71403c9656f19082e40e5a212efeeb29a640f05dc8bd7719e3d817272f659eed3ce29fb2492feee60bde5720c5c9cdd62e62fa295c3eef8dba618f0e01e408930ea3"
	inputsScriptHex += "ff"
	inputsScriptHex += "64222102c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71403c9656f19082e40e5a212efeeb29a640f05dc8bd7719e3d817272f659eed3ce29fb2492feee60bde5720c5c9cdd62e62fa295c3eef8dba618f0e01e408930ea3"
	inputsScriptHex += "ff"
	satisfiersHash := TxDoubleHash256(inputCountHex + inputsScriptHex)
	println("satisfiersHash:", satisfiersHash)

	txidem := "a87876c510a3823c041db9a04c6925014b8bd82a91862e721b4149d70d5a25c5"

	txIdHex := txidem + satisfiersHash

	txId := TxDoubleHash256(txIdHex)
	println("txId:", txId)
}
