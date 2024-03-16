package main

import (
	"fmt"
	"log"

	"github.com/kqns91/til/go/blockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Printf("w.PrivateKeyStr(): %v\n", w.PrivateKeyStr())
	fmt.Printf("w.PublicKeyStr(): %v\n", w.PublicKeyStr())
	fmt.Printf("w.BlockchainAddress(): %v\n", w.BlockchainAddress())
}
