package main

import (
	"fmt"
	"log"

	"github.com/Palaciodiego008/toonify"
)

func main() {
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice", "role": "admin"},
			{"id": 2, "name": "Bob", "role": "user"},
		},
	}

	encoded, err := toonify.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📤 TOON Output:")
	fmt.Println(encoded)

	// Decode back
	var decoded map[string]interface{}
	err = toonify.Decode(encoded, &decoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n📥 Decoded back:")
	fmt.Printf("%+v\n", decoded)
}
