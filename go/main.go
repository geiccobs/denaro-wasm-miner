//go:build js && wasm

package main

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/minio/sha256-simd"
	"log"
	"math"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

var (
	// config
	shareDifficulty = 0
	nodeUrl         = ""
	poolUrl         = ""
	serverUrl       = ""

	// stats
	shares      = 0
	minedBlocks = 0
)

func getTransactionsMerkleTree(transactions []string) string {

	var fullData []byte

	for _, transaction := range transactions {
		data, _ := hex.DecodeString(transaction)
		fullData = append(fullData, data...)
	}

	hash := sha256.New()
	hash.Write(fullData)

	return hex.EncodeToString(hash.Sum(nil))
}

func checkBlockIsValid(blockContent []byte, shareChunk string, chunk string, idifficulty int, charset string, hasDecimal bool) (bool, bool) {

	hash := sha256.New()
	hash.Write(blockContent)

	blockHash := hex.EncodeToString(hash.Sum(nil))

	if hasDecimal {
		return strings.HasPrefix(blockHash, shareChunk), strings.HasPrefix(blockHash, chunk) && strings.Contains(charset, string(blockHash[idifficulty]))
	} else {
		return strings.HasPrefix(blockHash, shareChunk), strings.HasPrefix(blockHash, chunk)
	}
}

func worker(start int, step int, res MiningInfoResult, address string) {

	var difficulty = res.Difficulty
	var idifficulty = int(difficulty)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from: %v\n", r)
			return
		}
	}()

	_, decimal := math.Modf(difficulty)

	lastBlock := res.LastBlock
	if lastBlock.Hash == "" {
		var num uint32 = 30_06_2005

		data := make([]byte, 32)
		binary.LittleEndian.PutUint32(data, num)

		lastBlock.Hash = hex.EncodeToString(data)
	}

	chunk := lastBlock.Hash[len(lastBlock.Hash)-idifficulty:]

	var shareChunk string

	if shareDifficulty > idifficulty {
		shareDifficulty = idifficulty
	}
	shareChunk = chunk[:shareDifficulty]

	charset := "0123456789abcdef"
	if decimal > 0 {
		count := math.Ceil(16 * (1 - decimal))
		charset = charset[:int(count)]
	}

	addressBytes := stringToBytes(address)
	t := float64(time.Now().UnixMicro()) / 1000000.0
	i := start
	a := time.Now().Unix()
	txs := res.PendingTransactionsHashes
	merkleTree := getTransactionsMerkleTree(txs)

	if start == 0 {
		log.Printf("Address: %s\n", address)
		log.Printf("Difficulty: %f\n", difficulty)
		log.Printf("Block number: %d\n", lastBlock.Id)
		log.Printf("Confirming %d transactions\n", len(txs))
	}

	var prefix []byte
	dataHash, _ := hex.DecodeString(lastBlock.Hash)
	prefix = append(prefix, dataHash...)
	prefix = append(prefix, addressBytes...)
	dataMerkleTree, _ := hex.DecodeString(merkleTree)
	prefix = append(prefix, dataMerkleTree...)
	dataA := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataA, uint32(a))
	prefix = append(prefix, dataA...)
	dataDifficulty := make([]byte, 2)
	binary.LittleEndian.PutUint16(dataDifficulty, uint16(difficulty*10))
	prefix = append(prefix, dataDifficulty...)

	if len(addressBytes) == 33 {
		data1 := make([]byte, 2, 2)
		binary.LittleEndian.PutUint16(data1, uint16(2))

		oldPrefix := prefix
		prefix = data1[:1]
		prefix = append(prefix, oldPrefix...)
	}

	for {
		var _hex []byte

		found := true
		check := 1000000 * step // 5000000 * step

	checkLoop:
		for {
			_hex = _hex[:0]
			_hex = append(_hex, prefix...)
			dataI := make([]byte, 4)
			binary.LittleEndian.PutUint32(dataI, uint32(i))
			_hex = append(_hex, dataI...)

			shareValid, blockValid := checkBlockIsValid(_hex, shareChunk, chunk, idifficulty, charset, decimal > 0)

			if shareValid {
				var txsInterface []interface{}
				for _, tx := range txs {
					txsInterface = append(txsInterface, tx)
				}

				js.Global().Call("expPostJSON", poolUrl+"share", map[string]interface{}{
					"block_content":    hex.EncodeToString(_hex),
					"txs":              txsInterface,
					"id":               lastBlock.Id + 1,
					"share_difficulty": difficulty,
				}, start+1)

				var response = js.Global().Get("response").Get(strconv.Itoa(start + 1))
				if response.Get("ok").Bool() {
					log.Printf("Worker n.%d: SHARE ACCEPTED", start+1)
					shares++
				} else {
					log.Printf("Worker n.%d: SHARE NOT ACCEPTED", start+1)
					log.Println(response.Get("error").String())
					return
				}
			}

			if blockValid {
				break checkLoop
			}

			i = i + step
			if (i-start)%check == 0 {
				elapsedTime := float64(time.Now().UnixMicro())/1000000.0 - t
				log.Printf("Worker %d: %dk hash/s", start+1, i/step/int(elapsedTime)/1000)

				js.Global().Call("expPostJSON", serverUrl+"setData?address="+address+"&worker_id="+strconv.Itoa(start), map[string]any{
					"hashrate":     i / step / int(elapsedTime) / 1000,
					"shares":       shares,
					"mined_blocks": minedBlocks,
					"last_update":  time.Now().Unix(),
				}, start+1)

				if elapsedTime > 90 {
					found = false
					break checkLoop
				}
			}
		}

		if found {
			var txsInterface []interface{}
			for _, tx := range txs {
				txsInterface = append(txsInterface, tx)
			}

			log.Println(hex.EncodeToString(_hex))

			js.Global().Call("expPostJSON", nodeUrl+"push_block", map[string]interface{}{
				"block_content": hex.EncodeToString(_hex),
				"txs":           txsInterface,
				"id":            lastBlock.Id + 1,
			}, start+1)

			var response = js.Global().Get("response").Get(strconv.Itoa(start + 1))
			if response.Get("ok").Bool() {
				log.Printf("Worker n.%d: BLOCK MINED", start+1)
				minedBlocks++
			} else {
				log.Println(response.Get("error").String())
			}
			return
		}
	}
}

func miner(_ js.Value, p []js.Value) any {
	miningAddress := p[0].String()

	nodeUrl = p[1].String()
	poolUrl = p[2].String()
	serverUrl = p[3].String()
	shareDifficulty, _ = strconv.Atoi(p[4].String())

	workerId := p[5].Int()
	workers := p[6].Int()

	for {
		log.Printf("Starting worker %d", workerId)

		js.Global().Call("expGetJSON", nodeUrl+"get_mining_info", workerId)

		miningInfo := js.Global().Call("expGetResponse", workerId).Get("result")
		lastBlock := miningInfo.Get("last_block")

		pendingTransactionHashes := make([]string, 0)

		for i := 0; i < miningInfo.Get("pending_transactions_hashes").Length(); i++ {
			pendingTransactionHashes = append(pendingTransactionHashes, miningInfo.Get("pending_transactions_hashes").Index(i).String())
		}

		reqP := MiningInfoResult{
			Difficulty: miningInfo.Get("difficulty").Float(),
			LastBlock: Block{
				Id:         int32(lastBlock.Get("id").Int()),
				Hash:       lastBlock.Get("hash").String(),
				Address:    lastBlock.Get("address").String(),
				Random:     int64(lastBlock.Get("random").Int()),
				Difficulty: lastBlock.Get("difficulty").Float(),
				Reward:     lastBlock.Get("reward").Float(),
				Timestamp:  int64(lastBlock.Get("timestamp").Int()),
			},
			PendingTransactionsHashes: pendingTransactionHashes,
			MerkleRoot:                miningInfo.Get("merkle_root").String(),
		}

		worker(workerId-1, workers, reqP, miningAddress)
	}
}

func main() {
	js.Global().Set("miner", js.FuncOf(miner))

	log.Println("WASM Go Initialized")
	<-make(chan bool)
}
