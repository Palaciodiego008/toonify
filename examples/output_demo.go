package main

import (
	"fmt"
	"log"

	"github.com/Palaciodiego008/toonify"
)

func main() {
	// 1. Simple object
	fmt.Println("=== SIMPLE OBJECT ===")
	simple := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"active": true,
	}
	
	output, err := toonify.Encode(simple)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// 2. Simple array
	fmt.Println("\n=== SIMPLE ARRAY ===")
	array := []interface{}{1, 2, 3, "hello", true}
	
	output, err = toonify.Encode(array)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// 3. Tabular array (users)
	fmt.Println("\n=== TABULAR ARRAY ===")
	users := []interface{}{
		map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
		map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
		map[string]interface{}{"id": 3, "name": "Charlie", "role": "guest"},
	}
	
	output, err = toonify.Encode(users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// 4. Complex structure
	fmt.Println("\n=== COMPLEX STRUCTURE ===")
	company := map[string]interface{}{
		"name": "TechCorp",
		"founded": 2020,
		"active": true,
		"employees": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice", "salary": 75000},
			map[string]interface{}{"id": 2, "name": "Bob", "salary": 65000},
		},
		"departments": []string{"Engineering", "Sales", "HR"},
		"metadata": map[string]interface{}{
			"version": "1.0",
			"updated": "2024-01-01",
		},
	}
	
	output, err = toonify.Encode(company)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// 5. Nested object
	fmt.Println("\n=== NESTED OBJECT ===")
	nested := map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"name": "Alice",
				"settings": map[string]interface{}{
					"theme": "dark",
					"notifications": true,
				},
			},
		},
	}
	
	output, err = toonify.Encode(nested)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
