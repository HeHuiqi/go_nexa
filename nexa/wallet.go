package nexa

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

// 由助记词的到十六进制的私钥和公钥，适用于BTC和ETH系列，HD Wallet

func MnemonicToPrivateKeyAndPublicKey(mnemonic string, path string, compress bool) (string, string) {

	fmt.Println("mnemonic:", mnemonic)

	seed := bip39.NewSeed(mnemonic, "")
	masterkey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	HardenedKeyStart := uint32(0x80000000) // 2^31
	derivekey0, _ := masterkey.Derive(HardenedKeyStart + 44)
	splitsStrs := strings.Split(path, "/")
	for i := 2; i < len(splitsStrs); i++ {
		indexStr := splitsStrs[i]
		// println("indexStr:", indexStr)
		isHard := strings.HasSuffix(indexStr, "'")
		index := uint32(0)
		if isHard {
			indexStr = indexStr[0 : len(indexStr)-1]
			tmp_index, _ := strconv.Atoi(indexStr)
			index = uint32(tmp_index) + HardenedKeyStart
		} else {
			tmp_index, _ := strconv.Atoi(indexStr)
			index = uint32(tmp_index)
		}
		derivekey0, _ = derivekey0.Derive(uint32(index))
		// println("indexStr:", indexStr)
	}

	pri, _ := derivekey0.ECPrivKey()
	// pri, _ := masterkey.ECPrivKey()

	pri_str := hex.EncodeToString(pri.Serialize())

	fmt.Println("pri_str:", pri_str)

	pub := pri.PubKey()
	var pubSerializeData []byte
	if compress {
		pubSerializeData = pub.SerializeCompressed()
	} else {
		pubSerializeData = pub.SerializeUncompressed()
	}

	pub_str := hex.EncodeToString(pubSerializeData)
	fmt.Println("pub_str:", pub_str)

	// 0x67C4341DD1fb25cf4e98624a40A76689CB61BfBd
	return pri_str, pub_str
}

func NexaMnemonicToPrivateKeyAndPublicKey(mnemonic string) (string, string) {
	path := "m/44'/29223'/0'/0/0" // receive address
	// path = "m/44'/29223'/0'/1/0"  // change address
	return MnemonicToPrivateKeyAndPublicKey(mnemonic, path, true)
}

func NexaPrivateKeyToPublicKeyAndAddress(privateKeyHex string) (string, string) {
	priBytes, _ := hex.DecodeString(privateKeyHex)
	_, publicKey := btcec.PrivKeyFromBytes(priBytes)
	pubkeyData := publicKey.SerializeCompressed()
	publickHex := hex.EncodeToString(pubkeyData)
	nexaAddress := NexaPublickeyToAddress(publickHex)
	return publickHex, nexaAddress
}

func NexaPublickeyToAddress(publickHex string) string {
	prefix := "nexa"
	pubkeyData, _ := hex.DecodeString(publickHex)
	publen := uint8(len(pubkeyData))
	pubkeyData = append([]byte{publen}, pubkeyData...)
	// all(len) + op_o + op_1 + hash160(len(pubkey)+pubkey )
	// 17 00 51 14 + hash160(pubkey)
	format, _ := hex.DecodeString("17005114")
	pubkeyDataH160 := append(format, btcutil.Hash160(pubkeyData)...)
	// pubkeyDataH160, _ = hex.DecodeString("170051147c948439fe65e511b7737fe69be921d6531412be")
	// println("pubkeyDataH160:", hex.EncodeToString(pubkeyDataH160))
	nexaAddress := CheckEncodeCashAddress(pubkeyDataH160, prefix, SCRIPT_TEMPLATE_TYPE)
	nexaAddress = prefix + ":" + nexaAddress
	println("nexaAddress:", nexaAddress)
	return nexaAddress
}
