package main

import (
	"log"

	"gomasterforge/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
