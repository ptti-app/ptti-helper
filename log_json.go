package ptti

import (
	"encoding/json"
	"log"
)

// convert to JSON and log out
func Log(T any, customMsg ...string) {
	json, err := json.MarshalIndent(T, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	if customMsg != nil {
		log.Printf("%s: %v\n", customMsg, string(json))
	} else {
		log.Printf("Results: %v\n", string(json))
	}
}
