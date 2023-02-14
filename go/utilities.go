package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
)

func stringToBytes(text string) []byte {

	var data []byte

	data, err := hex.DecodeString(text)
	if err != nil {
		data = base58.Decode(text)
	}

	return data
}
