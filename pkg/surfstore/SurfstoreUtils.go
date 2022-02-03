package surfstore

import "log"

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) {
	var blockStoreAddr string
	if err := client.GetBlockStoreAddr(&blockStoreAddr); err != nil {
		log.Fatal(err)
	}

	blockHash := "1234"
	var block Block
	if err := client.GetBlock(blockHash, blockStoreAddr, &block); err != nil {
		log.Fatal(err)
	}

	log.Print(block.String())
}
