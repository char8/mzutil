package main

import (
	"context"
	"fmt"
	"github.com/char8/mzutil/auth"
	"github.com/char8/mzutil/client"
	"io/ioutil"
)

func main() {
	ks := client.KeychainConfigStore{}

	//ctx := context.Background()

	a := auth.NewAuthenticator(ks, 10035)
	//err := a.Login()

	cli := a.NewClient(context.Background())
	resp, err := cli.Get("https://api.monzo.com/ping/whoami")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\n\n", string(body))
	fmt.Printf("%v\n", err)
}
