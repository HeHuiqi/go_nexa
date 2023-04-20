package nexa

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
)

func VarIntToHex(val uint64) string {
	w := bytes.NewBuffer(make([]byte, 0, 8))
	wire.WriteVarInt(w, uint32(0), val)
	hex := hex.EncodeToString(w.Bytes())
	return hex
}
