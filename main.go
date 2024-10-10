package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
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
	accountList := []ton.Address{tongo.MustParseAddress("0:A560DC2F757589A7152835573D23117910207DDEDEFA717DBC82D8CF932E31B2"),
		tongo.MustParseAddress("-1:CFB0B11350E4052FB60ED14CBF2FEDB265D7AC9EEB4C6BE0EEEE60AA2CE5DB3E"),
		tongo.MustParseAddress("-1:5555555555555555555555555555555555555555555555555555555555555555"),
		tongo.MustParseAddress("-1:0000000000000000000000000000000000000000000000000000000000000000"),
		tongo.MustParseAddress("0:A560DC2F757589A7152835573D23117910207DDEDEFA717DBC82D8CF932E31B2")}
	testTime := 10000 * time.Second

	req_count := 10
	total_req_count := 0
	for j := 0; j < req_count; j++ {
		call := func() {
			count := 0
			totalElapsed := time.Duration(0)
			for {
				count++
				total_req_count++
				start := time.Now()
				_, err := tongoClient.GetAccountState(context.Background(), accountList[j%len(accountList)].ID)
				if err != nil {
					fmt.Printf("Get account state error: %v", err)
					return
				}
				totalElapsed += time.Since(start)
				if count%100 == 0 {
					fmt.Printf("%v, GetAccountState avg cost: %v ms，totalElapsed:%v, count:%v, total_req_count:%v \n", j, int(totalElapsed.Milliseconds())/count, totalElapsed, count, total_req_count)
					count = 0
					totalElapsed = 0
				}
			}
		}

		go call()
	}
	time.Sleep(testTime)
}
