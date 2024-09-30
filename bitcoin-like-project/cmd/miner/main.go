package main

import (
	"log"

	"github.com/TanmayDhobale/bitcoin-like-project/internal/blockchain"
	"github.com/TanmayDhobale/bitcoin-like-project/internal/miner"
)

func main() {
	
	bc := blockchain.NewBlockchain("miner_address")
	minerNode := miner.NewMinerNode(bc)

	err := minerNode.ConnectToCentralServer("ws://localhost:8080/ws")
	if err != nil {
		log.Fatal("Error connecting to central server:", err)
	}

	log.Println("Miner connected to central server. Starting mining...")
	minerNode.StartMining()
}