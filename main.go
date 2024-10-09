package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
)

func main() {
	// config, err := config.ParseConfigFile("./config.json")
	// tongoClient, err := liteapi.NewClientWithDefaultTestnet()
	// tongoClient, err := liteapi.NewClient(liteapi.WithConfigurationFile(*config))
	tongoClient, err := liteapi.NewClient(liteapi.WithLiteServers([]config.LiteServer{{Host: "127.0.0.1:4443", Key: "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY="}}))
	if err != nil {
		fmt.Printf("Unable to create tongo client: %v", err)
	}
	// accountId := tongo.MustParseAccountID("0:E2D41ED396A9F1BA03839D63C5650FAFC6FCFB574FD03F2E67D6555B61A3ACD9")
	account := tongo.MustParseAddress("0:45CF4642228CB12F514BB3A79C389A7BAA7A10A92E8ECF4B915D246B7C2708CB")

	var wg sync.WaitGroup
	req_count := 1000
	for i := 0; i < 1000/req_count; i++ {
		for j := 0; j < req_count; j++ {
			call := func() {
				defer wg.Done()

				start := time.Now()
				_, err := tongoClient.GetAccountState(context.Background(), account.ID)
				elapsed := time.Now().Sub(start)
				fmt.Printf("GetAccountState cost: %s\n", elapsed)
				if err != nil {
					fmt.Printf("Get account state error: %v", err)
				}
			}

			wg.Add(1)
			go call()
		}

		wg.Wait()
		fmt.Println("------------")
		// fmt.Printf("Account status: %v\nBalance: %v\n", state.Account.Status(), state.Account.Account.Storage.Balance.Grams)
	}
}
