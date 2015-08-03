package main

import (
	"flag"
	"fmt"

	"github.com/jeevatkm/urlite/util"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "admin", "Password for administrator accout")
	bcost := flag.Int("bcryptCost", bcrypt.DefaultCost, `bcrypt cost value. default value is good enough i.e 10
		Kindly make sure this cost factor is same as one in the
		'/etc/urlite/urlite.conf'`)

	flag.Parse()

	hashPass := util.HashPassword(*password)

	fmt.Println("\n\nurlite password generate tool")
	fmt.Println("-----------------------------")
	fmt.Printf("Password string	: %s\nBcrypt cost	: %d", *password, *bcost)
	fmt.Printf("\n\nBcrypt password hash: %s\n\n", hashPass)
}
