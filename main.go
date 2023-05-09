package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slices"
)

type BlockReturnStruct struct {
	Jsonrpc any               `json:"jsonrpc"`
	Id      any               `json:"id"`
	Result  BlockResultStruct `json:"result,omitempty"`
	Error   ErrorStruct       `json:"error,omitempty"`
}

type BlockResultStruct struct {
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

type TxReturnStruct struct {
	Jsonrpc any            `json:"jsonrpc"`
	Id      any            `json:"id"`
	Result  TxResultStruct `json:"result,omitempty"`
	Error   ErrorStruct    `json:"error,omitempty"`
}

type TxResultStruct struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	ChainId          string `json:"chainId"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
	V                string `json:"v"`
	Value            string `json:"value"`
}

type ErrorStruct struct {
	Code    any `json:"code"`
	Message any `json:"message"`
}

var RpcUrl = "http://192.168.219.105:8545"
var FileName = "./Ethlistener.data"

var DB *sql.DB
var address []string

func main() {
	DB, _ = sql.Open("mysql", "root:1234qwer@tcp(192.168.219.107:13306)/Inae")

	//gocron.Every(1).Minutes().Do(cronJob)
	//<-gocron.Start()

	cronJob()
}

func cronJob() {
	address, _ = getAddress()

	latest := getLatest()
	height := getBlockCount()

	latestInt := hex2int(latest)
	heightInt := hex2int(height)

	for latestInt < heightInt {
		latestInt += 1
		getBlock(int2hex(latestInt))
	}

	saveLatest(height)
}

func hex2int(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return result
}

func int2hex(intVal uint64) string {
	return fmt.Sprintf("0x%x", intVal)
}

func getAddress() ([]string, error) {
	var result []string

	query := `SELECT address FROM eth_address`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := ""

		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func saveLatest(height string) {
	err := ioutil.WriteFile(FileName, []byte(height), 0)
	if err != nil {
		log.Fatalf("Write File: %v", err)
	}
}

func getLatest() string {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Fatalf("Read File: %v", err)
	}

	return string(data)
}

func getBlockCount() string {
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"id":      0,
		"params":  []interface{}{},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}

	resp, err := http.Post(RpcUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return fmt.Sprintf("%v", result["result"])
}

func getBlock(blocknumber string) {
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"id":      0,
		"params":  []interface{}{blocknumber, false},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}

	resp, err := http.Post(RpcUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}

	result := BlockReturnStruct{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	dbtx, err := DB.Begin()
	if err != nil {
		log.Fatalf("DB Transaction: %v", err)
	}

	defer dbtx.Rollback()

	for _, tx := range result.Result.Transactions {
		getTransaction(tx, dbtx)
	}

	err = dbtx.Commit()
	if err != nil {
		log.Fatalf("DB Commit: %v", err)
	}
}

func getTransaction(hash string, dbtx *sql.Tx) {

	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"id":      0,
		"params":  []interface{}{hash},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}

	resp, err := http.Post(RpcUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}

	result := TxReturnStruct{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	query := `INSERT INTO eth_history(block_number, block_hash, tx_hash, from_address, to_address, value, create_dt) VALUES (?, ?, ?, ?, ?, ?, NOW())`

	idx := slices.IndexFunc(address, func(c string) bool { return c == result.Result.From || c == result.Result.To })
	if idx != -1 {
		_, err = dbtx.Exec(query, hex2int(result.Result.BlockNumber), result.Result.BlockHash, result.Result.Hash, result.Result.From, result.Result.To, hex2int(result.Result.Value))
		if err != nil {
			log.Fatalf("DB Insert: %v", err)
		}
	}
}
