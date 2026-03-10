package main

import (
	"fmt"
	"github.com/Avi-88/lusus/store"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func main() {

	s := store.NewLususStore("demoDB")
	reader := bufio.NewReader(os.Stdin)

	for {
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		args := strings.Split(command, " ")
		
		if len(args) == 0 {
			continue
		}
	
		operation := strings.ToUpper(args[0])
		switch operation {
		case "SET":
			if len(args) < 3 || len(args) > 4 {
				fmt.Println("Usage SET <key> <value> [optional]<ttl_sec>")
				continue
			} 
			var err error
			if len(args) == 4 {
				ttl, parseErr := strconv.Atoi(args[3])
				if parseErr != nil {
					fmt.Println("Invalid TTL")
					continue
				}
				err = s.Set(args[1], args[2], ttl)
			}else{
				err = s.Set(args[1], args[2])
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Key stored")	

		case "GET":
			if len(args) != 2 {
				fmt.Println("Usage GET <key>")
				continue
			}
			val, exists := s.Get(args[1])
			if !exists {
				fmt.Println("Key not found")
				continue
			}
			fmt.Println("Value:", val)
		
		case "DEL":
			if len(args) != 2 {
				fmt.Println("Usage DEL <key>")
				continue
			}
			err := s.Delete(args[1])
			if err != nil {
				fmt.Println("Error:", err)
			}else{
				fmt.Println("Key deleted")
			}

		case "EXPIRE":
			if len(args) != 3 {
				fmt.Println("Usage EXPIRE <key> <ttl_sec>")
				continue
			}
			ttl, parseErr := strconv.Atoi(args[2])
			if parseErr != nil {
				fmt.Println("Invalid TTL")
				continue
			}
			ok := s.Expire(args[1], ttl)
			if !ok {
				fmt.Println("Key not found or already expired")
				continue
			}
			fmt.Println("Key expiry set")

		case "TTL":
			if len(args) != 2 {
				fmt.Println("Usage TTL <key>")
				continue
			}
			out := s.TTL(args[1])
			switch out {
			case -2:
				fmt.Println("Key not found")
			case -1:
				fmt.Println("Key not set to expire")
			default:
				fmt.Printf("Key expires in %d secs\n", out)
			}
		case "EXIT":
			fmt.Println("Exiting CLI")
			return

		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}