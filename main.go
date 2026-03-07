package main

import (
	"fmt"
	"github.com/Avi-88/lusus/store"
)

func main() {

	s := store.NewLususStore()
	s.Set("Name", "Avi")

	val, ok := s.Get("Name")
	if !ok {
		fmt.Println("Key not found")
	}
	fmt.Printf("Created a new store and stored the value %s",val)
}