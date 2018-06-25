package main

import (
	// "context"
	// "fmt"
	// "io/ioutil"
	// "log"
	//
	// "github.com/char8/mzutil/config"
	// "github.com/char8/mzutil/monzo"
	"github.com/char8/mzutil/cmd"
)

func main() {
	cmd.Execute()
	// ks := config.NewKeychainConfigStore("mzutil")
	//
	// //ctx := context.Background()
	//
	// a, err := monzo.NewAuthenticator(ks)
	// if err != nil {
	// 	log.Fatalf("Could not create authenticator: %v", err)
	// }
	// err = a.Login()
	//
	// cli := a.NewClient(context.Background())
	// resp, err := cli.Get("https://api.monzo.com/ping/whoami")
	//
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }
	//
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("%v\n\n", string(body))
	// fmt.Printf("%v\n", err)
}
