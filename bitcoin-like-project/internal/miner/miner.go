package miner

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/TanmayDhobale/bitcoin-like-project/internal/blockchain"
)

type MinerNode struct {
	blockchain *blockchain.Blockchain
	conn       *websocket.Conn
}

func NewMinerNode(bc *blockchain.Blockchain) *MinerNode {
	return &MinerNode{
		blockchain: bc,
	}
}

func (m *MinerNode) ConnectToCentralServer(url string) error {
	var err error
	m.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	go m.readMessages()
	return nil
}

func (m *MinerNode) readMessages() {
	for {
		_, message, err := m.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		log.Printf("Received message: %s", message)
	}
}

func (m *MinerNode) StartMining() {
	for {
		newBlock := m.mineBlock()
		m.broadcastBlock(newBlock)
		time.Sleep(10 * time.Second)
	}
}

func (m *MinerNode) mineBlock() *blockchain.Block {
	prevBlock := m.blockchain.GetBlocks()[len(m.blockchain.GetBlocks())-1]
	coinbaseTx := blockchain.NewCoinbaseTX("miner_address", "")
	newBlock := blockchain.NewBlock([]*blockchain.Transaction{coinbaseTx}, prevBlock.Hash, m.blockchain.AdjustDifficulty())
	err := m.blockchain.AddBlock([]*blockchain.Transaction{coinbaseTx})
	if err != nil {
		log.Printf("Error adding block: %v", err)
	}
	return newBlock
}

func (m *MinerNode) broadcastBlock(block *blockchain.Block) {
	blockData := struct {
		Type string             `json:"type"`
		Data *blockchain.Block `json:"data"`
	}{
		Type: "new_block",
		Data: block,
	}
	err := m.conn.WriteJSON(blockData)
	if err != nil {
		log.Println("Error broadcasting block:", err)
	}
}