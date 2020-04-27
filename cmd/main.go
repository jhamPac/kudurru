package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jhampac/kudurru"
)

func main() {
	s := kudurru.New()

	fmt.Printf("Kudurru serving on port:%v", os.Getenv("PORT"))
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
