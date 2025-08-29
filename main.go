package main

import (
	"fmt"
	"github.com/spencerreeves/gophercon-challenge-2025/freeway"
	"log"
)

func main() {
	f1, err := freeway.PortKnockV2("freeway.white-rabbit.dev:30152")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Flag: %v", f1)
}
