package nexa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

func Hash256(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

// mac returns an HMAC of the given key and message.
func Mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// key value=message
func Sha256hmac(k []byte, m []byte) []byte {
	alg := sha256.New
	v := Mac(alg, k, m)
	return v
}

func BigFromHex(s string) *big.Int {
	if s == "" {
		return big.NewInt(0)
	}
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex in source file: " + s)
	}
	return r
}

func BigFromBytes(bytes []byte) *big.Int {
	return big.NewInt(0).SetBytes(bytes)
}

func NonceFunctionRFC6979(privkey []byte, msgbuf []byte) *big.Int {
	V, _ := hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
	K, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	additionalData := []byte{'S', 'c', 'h', 'n', 'o', 'r', 'r', '+', 'S', 'H', 'A', '2', '5', '6', ' ', ' '}

	space := []byte{}
	blob := append(privkey, msgbuf...)
	blob = append(blob, space...)
	blob = append(blob, additionalData...)
	// println("blob:", hex.EncodeToString(blob))
	// V+b'\x00'+blob
	K = Sha256hmac(K, append(append(V, 0x00), blob...))
	V = Sha256hmac(K, V)
	// println("K:", hex.EncodeToString(K))
	// println("V:", hex.EncodeToString(V))

	K = Sha256hmac(K, append(append(V, 0x01), blob...))
	V = Sha256hmac(K, V)
	// println("K:", hex.EncodeToString(K))
	// println("V:", hex.EncodeToString(V))

	k := big.NewInt(0)

	for {
		V = Sha256hmac(K, V)
		T := V
		// println("V:", hex.EncodeToString(V))
		// println("T:", hex.EncodeToString(T))
		k = BigFromBytes(T)
		zero := big.NewInt(0)
		if k.Cmp(zero) > 0 && k.Cmp(GetN()) < 0 {
			break
		}
		K = Sha256hmac(K, append(V, 0x00))
		K = Sha256hmac(K, V)
	}
	// println("k", hex.EncodeToString(k.Bytes()))
	// println("k", k.String())

	return k
}

/*
Schnorr Sig可以与ECDSA使用同一个椭圆曲线：secp256k1 curve，升级起来的改动非常小。
原理
我们定义几个变量：

G：椭圆曲线。
m：待签名的数据，通常是一个32字节的哈希值。
x：私钥。P = xG，P为x对应的公钥。
H()：哈希函数。
示例：写法H(m || R || P)可理解为：将m, R, P三个字段拼接在一起然后再做哈希运算。
生成签名
签名者已知的是：G-椭圆曲线, H()-哈希函数，m-待签名消息, x-私钥。

选择一个随机数k, 令 R = kG
令 s = k + H(m || R || P)*x
那么，公钥P对消息m的签名就是：(R, s)，这一对值即为Schnorr签名

*/

// 返回签名和公钥
func Signature(pri []byte, m []byte) (string, string) {
	// P = G * k   P公钥 G椭圆曲线 k私钥
	n := GetN()
	// P = G*d // 就是由私钥匙获取公钥
	// R = G*k // 同上

	dPrivKey, P := btcec.PrivKeyFromBytes(pri)
	dPubkeyBytes := P.SerializeCompressed()

	println("dPrivKey:", dPrivKey.Key.String())

	k := NonceFunctionRFC6979(pri, m)
	kPrivKey, R := btcec.PrivKeyFromBytes(k.Bytes())
	println("kPrivKey:", kPrivKey.Key.String())

	p := GetP()
	if big.Jacobi(R.Y(), p) == -1 {
		println("k-Sub-Sub-Sub")
		k = k.Sub(n, k)
	}

	r := R.X()
	// hash(r + p + m)
	hashMessage := append(r.Bytes(), dPubkeyBytes...)
	hashMessage = append(hashMessage, m...)
	e0 := BigFromBytes(Hash256(hashMessage))

	// s = (e0*pri + k)
	priNum := big.NewInt(0).SetBytes(pri)
	num := e0.Mul(e0, priNum)
	num = num.Add(num, k)

	s := num.Mod(num, n)

	hexR := hex.EncodeToString(r.Bytes())
	hexS := hex.EncodeToString(s.Bytes())
	// println("r:", hexR)
	// println("s:", hexS)

	sign := hexR + hexS

	// println("sign:", sign)
	return sign, hex.EncodeToString(dPubkeyBytes)

}
