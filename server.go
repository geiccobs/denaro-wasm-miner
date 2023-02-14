package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	port = ":3010"
	dir  = "."
)

/*
data format is:

	address: {
		worker: {
			hashrate: int
			shares: int
			mined_blocks: int
			last_update: int64
		}
	}
*/
var data = make(map[string]map[string]map[string]any) // it's almost funny how much I hate this

func setData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if req.Method == "OPTIONS" {
		return
	}

	q := req.URL.Query()
	address := q.Get("address")
	workerId := q.Get("worker_id")

	if len(address) < 32 {
		return
	}

	// decode JSON
	var tmp map[string]any
	if err := json.NewDecoder(req.Body).Decode(&tmp); err != nil {
		panic(err)
	}

	delete(tmp, "address")

	if _, ok := data[address]; !ok {
		data[address] = make(map[string]map[string]any)
	}
	data[address][workerId] = tmp
	data[address][workerId]["last_update"] = time.Now().Unix()

	// return ok status
	if err := json.NewEncoder(w).Encode(map[string]any{"result": "ok"}); err != nil {
		panic(err)
	}
}

func getData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	address := req.URL.Query().Get("address")

	// check if data exists
	if len(address) == 0 {
		// return data
		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	} else {
		// return data
		if err := json.NewEncoder(w).Encode(getTotal(address)); err != nil {
			panic(err)
		}
	}
}

func getTotal(address string) map[string]any {
	if _, ok := data[address]; !ok {
		return map[string]any{"hashrate": 0, "shares": 0, "mined_blocks": 0, "last_update": 0}
	}

	total := make(map[string]any)
	for workerId, worker := range data[address] {
		for key, value := range worker {
			if key == "last_update" {
				// more than 60 seconds since last update
				if time.Now().Unix()-value.(int64) > 60 {
					delete(data[address], workerId)
				}
				continue
			}

			if _, ok := total[key]; !ok {
				total[key] = 0
			}
			total[key] = total[key].(int) + int(value.(float64))
		}
	}
	return total
}

func proxyRequest(w http.ResponseWriter, req *http.Request) {
	// some websites don't allow cross-origin requests, that's when this function comes in handy
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if req.Method == "OPTIONS" {
		return
	}

	// execute request based on method
	var fwRequest *http.Response
	var err error

	switch req.Method {
	case "GET":
		if fwRequest, err = http.Get(req.URL.Query().Get("url")); err != nil {
			panic(err)
		}
	case "POST":
		if fwRequest, err = http.Post(req.URL.Query().Get("url"), "application/json", req.Body); err != nil {
			panic(err)
		}
	}

	// return body as data
	if _, err = io.Copy(w, fwRequest.Body); err != nil {
		panic(err)
	}
}

func main() {
	println("Listening on port", port)

	http.HandleFunc("/setData", setData)
	http.HandleFunc("/getData", getData)
	http.HandleFunc("/proxy/", proxyRequest)

	http.Handle("/", http.FileServer(http.Dir(dir)))

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
