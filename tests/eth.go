package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)
	func main() {
	client, err := ethclient.Dial("HTTP://127.0.0.1:8545")
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("we have a connection")
	chainID,err := client.NetworkID(context.Background())
	if err != nil{
		fmt.Println("error occured")
	}
	fmt.Println(chainID.Int64())


	// _ = client // we'll use this in the upcoming sections
	}