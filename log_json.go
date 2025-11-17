package helper

import (
	"encoding/json"
	"log"
)

// convert to JSON and log out
func Log(T any) {
	json, err := json.MarshalIndent(T, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Results: %v\n", string(json))
}
