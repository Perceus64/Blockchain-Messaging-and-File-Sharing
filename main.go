package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Description  string
	SenderHash   string
	ReceiverHash string
	FileHash     string
	PrevHash     string
	Hash         string
}

var Blockchain []Block

// EmailToHash maps email addresses to their corresponding hash
var EmailToHash map[string]string

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Description + block.SenderHash + block.ReceiverHash + block.FileHash + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, description, senderHash, receiverHash, fileHash string) Block {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Description = description
	newBlock.SenderHash = senderHash
	newBlock.ReceiverHash = receiverHash
	newBlock.FileHash = fileHash
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {
	genesisBlock := Block{}
	genesisBlock = Block{0, time.Now().String(), "Genesis Block", "", "", "", "", calculateHash(genesisBlock)}
	Blockchain = append(Blockchain, genesisBlock)

	EmailToHash = make(map[string]string)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Send Message")
		fmt.Println("2. Send File")
		fmt.Println("3. Send Message and File")
		fmt.Println("4. Print Blockchain and Quit")

		optionStr, _ := reader.ReadString('\n')
		optionStr = strings.TrimSpace(optionStr)
		option := strings.ToLower(optionStr)

		switch option {
		case "1":
			sendMessage(reader)
		case "2":
			sendFile(reader)
		case "3":
			sendMessageAndFile(reader)
		case "4":
			printBlockchain()
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please choose again.")
		}
	}
}

func printBlockchain() {
	fmt.Println("Blockchain Contents:")
	for _, block := range Blockchain {
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Description:", block.Description)
		fmt.Println("SenderHash:", block.SenderHash)
		fmt.Println("ReceiverHash:", block.ReceiverHash)
		fmt.Println("FileHash:", block.FileHash)
		fmt.Println("PrevHash:", block.PrevHash)
		fmt.Println("Hash:", block.Hash)
		fmt.Println()
	}
}

func sendMessage(reader *bufio.Reader) {
	fmt.Println("Enter action description:")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Println("Enter sender's email:")
	senderEmail, _ := reader.ReadString('\n')
	senderEmail = strings.TrimSpace(senderEmail)

	fmt.Println("Enter receiver's email:")
	receiverEmail, _ := reader.ReadString('\n')
	receiverEmail = strings.TrimSpace(receiverEmail)

	fmt.Println("Enter message:")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	senderHash := fmt.Sprintf("%x", sha256.Sum256([]byte(senderEmail)))
	receiverHash := fmt.Sprintf("%x", sha256.Sum256([]byte(receiverEmail)))

	EmailToHash[senderEmail] = senderHash
	EmailToHash[receiverEmail] = receiverHash

	fmt.Println("Message sent successfully!")

	previousBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(previousBlock, description, senderHash, receiverHash, message)
	Blockchain = append(Blockchain, newBlock)

	fmt.Println("Block added to the blockchain!")
}

func sendFile(reader *bufio.Reader) {
	fmt.Println("Enter action description:")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Println("Enter sender's email:")
	senderEmail, _ := reader.ReadString('\n')
	senderEmail =
		strings.TrimSpace(senderEmail)

	fmt.Println("Enter receiver's email:")
	receiverEmail, _ := reader.ReadString('\n')
	receiverEmail = strings.TrimSpace(receiverEmail)

	fmt.Println("Enter file path (e.g., /path/to/file.txt):")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileSize := fileInfo.Size()

	fileDetails := fmt.Sprintf("File Name: %s\nFile Size: %d bytes", fileInfo.Name(), fileSize)

	senderHash := fmt.Sprintf("%x", sha256.Sum256([]byte(senderEmail)))
	receiverHash := fmt.Sprintf("%x", sha256.Sum256([]byte(receiverEmail)))

	EmailToHash[senderEmail] = senderHash

	fmt.Println("File sent successfully!")

	previousBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(previousBlock, description, senderHash, receiverHash, fileDetails)
	Blockchain = append(Blockchain, newBlock)

	fmt.Println("Block added to the blockchain!")
}

func sendMessageAndFile(reader *bufio.Reader) {
	sendMessage(reader)
	sendFile(reader)
}
