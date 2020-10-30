package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/dmc"
)

func main() {
	scope := flag.String("scope", "", "Scope definition to be used for the machine credential")
	flag.Parse()
	if *scope == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	rpc := dmc.NewLRPC2()
	token, err := rpc.GetToken(*scope)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(token)
}
