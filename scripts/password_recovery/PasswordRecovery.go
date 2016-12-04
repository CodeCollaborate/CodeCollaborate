package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	password = flag.String("password", "password", "the unhashed password you want the hash of")
)

func main() {
	flag.Parse()

	args := os.Args[1:]
	if len(args) > 0 && *password == "password" {
		// didn't use flag
		inferred := args[len(args)-1]
		password = &inferred
	}

	fmt.Printf("generating hash for password %s: \n", string(*password))

	hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR: problem making hash")
		fmt.Println(err)
	}
	fmt.Println(string(hashed))
}
