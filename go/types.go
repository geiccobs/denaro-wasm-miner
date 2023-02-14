package main

type MiningInfoResult struct {
	Difficulty                float64     `json:"difficulty"`
	LastBlock                 Block       `json:"last_block"`
	PendingTransactions       interface{} `json:"pending_transactions"`
	PendingTransactionsHashes []string    `json:"pending_transactions_hashes"`
	MerkleRoot                string      `json:"merkle_root"`
}

type Block struct {
	Id         int32   `json:"id"`
	Hash       string  `json:"hash"`
	Address    string  `json:"address"`
	Random     int64   `json:"random"`
	Difficulty float64 `json:"difficulty"`
	Reward     float64 `json:"reward"`
	Timestamp  int64   `json:"timestamp"`
}
