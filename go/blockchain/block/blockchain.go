package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kqns91/til/go/blockchain/utils"
)

const (
	// ハッシュ値の先頭に必要なゼロの数。(マイニングの難易度)
	MINING_DIFFICULTY = 3
	// マイニング報酬の送信者。
	MINING_SENDER = "THE BLOCKCHAIN"
	// マイニング報酬。
	MINING_REWARD    = 1.0
	MINING_TIMER_SEC = 20

	BLCOKCHAIN_PORT_RANGE_START       = 5001
	BLCOKCHAIN_PORT_RANGE_END         = 5004
	NEIGHBOR_IP_RANGE_START           = 0
	NEIGHBOR_IP_RANGE_END             = 1
	BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC = 20
)

// ブロックの構造体。
type Block struct {
	// ブロックが生成された時刻。
	timestamp int64
	// ブロックのナンス。
	nonce int
	// 1つ前のブロックのハッシュ値。
	previousHash [32]byte
	// ブロックで処理するトランザクションのリスト。
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) Print() {
	fmt.Printf("timestamp     : %d\n", b.timestamp)
	fmt.Printf("nonce         : %d\n", b.nonce)
	fmt.Printf("previousHash  : %x\n", b.previousHash)
	for _, v := range b.transactions {
		v.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

// ブロックチェーンの構造体。
type Blockchain struct {
	// トランザクションプール。
	transactionPool []*Transaction
	// ブロックのリスト。
	chain []*Block
	// マイニング報酬の受信者のブロックチェーンアドレス。
	blockchainAddress string
	port              uint16
	mux               sync.Mutex

	neighbors    []string
	muxNeighbors sync.Mutex
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := &Blockchain{blockchainAddress: blockchainAddress, port: port}
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
}

func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(), bc.port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLCOKCHAIN_PORT_RANGE_START, BLCOKCHAIN_PORT_RANGE_END,
	)
	log.Printf("%v", bc.neighbors)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC*time.Second, bc.StartSyncNeighbors)
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) CreateTransaction(
	senderBlockchainAddress, recipientBlockchainAddress string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature,
) bool {
	isTransacted := bc.AddTransaction(senderBlockchainAddress, recipientBlockchainAddress, value, senderPublicKey, s)

	if isTransacted {
		for _, n := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				SenderBlockchainAddress:    &senderBlockchainAddress,
				RecipientBlockchainAddress: &recipientBlockchainAddress,
				SenderPublicKey:            &publicKeyStr,
				Value:                      &value,
				Signature:                  &signatureStr,
			}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest(http.MethodPut, endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("%v", resp)
		}
	}

	return isTransacted
}

// トランザクションを追加する。
func (bc *Blockchain) AddTransaction(
	senderBlockchainAddress, recipientBlockchainAddress string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature,
) bool {
	t := NewTransaction(senderBlockchainAddress, recipientBlockchainAddress, value)

	if senderBlockchainAddress == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(senderBlockchainAddress) < value {
		// 	log.Println("ERROR: Not enough balance in a wallet")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	log.Println("ERROR: Verify Transaction")

	return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublisKey *ecdsa.PublicKey,
	s *utils.Signature,
	t *Transaction,
) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublisKey, h[:], s.R, s.S)
}

// トランザクションプールをコピーする。
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value,
			),
		)
	}
	return transactions
}

// ハッシュ値の検証。
func (bc *Blockchain) ValidProof(
	nonce int,
	previousHash [32]byte,
	transactions []*Transaction,
	difficulty int,
) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

// 有効なハッシュ値が見つかるまで、ナンスをインクリメントする。
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// マイニングを行う。
func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.transactionPool) == 0 {
		return false
	}

	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(MINING_TIMER_SEC*time.Second, bc.StartMining)
}

// ブロックチェーンアドレスの残高を計算する。
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// トランザクションの構造体。
type Transaction struct {
	// 送信者のブロックチェーンアドレス。
	senderBlockchainAddress string
	// 受信者のブロックチェーンアドレス。
	recipientBlockchainAddress string
	// 送信する量。
	value float32
}

func NewTransaction(senderBlockchainAddress, recipientBlockchainAddress string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    senderBlockchainAddress,
		recipientBlockchainAddress: recipientBlockchainAddress,
		value:                      value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address    : %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address : %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                        : %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (t *TransactionRequest) Validate() bool {
	return t.SenderBlockchainAddress != nil &&
		t.RecipientBlockchainAddress != nil &&
		t.SenderPublicKey != nil &&
		t.Value != nil &&
		t.Signature != nil
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}
