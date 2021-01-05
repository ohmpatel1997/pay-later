package main

import (
	"bufio"
	"fmt"
	"os"
	"pay-later/integration/log"
	"pay-later/service/command"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {

		l := log.NewLogger()
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		text = strings.Replace(text, "\n", "", -1)

		srv, err := command.NewCommand(text)
		if err != nil {
			fmt.Println(err)
		}

		srv.Execute(l)
	}
}

//new user user1 u1@users.com 300
//new user user2 u2@users.com 400
//new user user3 u3@users.com 500
//new merchant m1 m1@merchants.com 0.5%
//new merchant m2 m2@merchants.com 1.5%
//new merchant m3 m3@merchants.com 1.25%
//new txn user2 m1 500
//new txn user1 m2 300
//new txn user1 m3 10
//report users-at-credit-limit
//new txn user3 m3 200
//new txn user3 m3 300
//report users-at-credit-limit
//report discount m3
//payback user3 400
//report total-dues
