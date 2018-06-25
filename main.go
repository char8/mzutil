package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/char8/mzutil/config"
	"github.com/char8/mzutil/monzo"
)

func main() {
	ks := config.NewKeychainConfigStore("mzutil")

	//ctx := context.Background()

	a := monzo.NewAuthenticator(ks, 10035)
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
