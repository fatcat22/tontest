package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
)

type ReqResult struct {
	elapsed time.Duration
	err     error
}

func main() {
	// config, err := config.ParseConfigFile("./config.json")
	// tongoClient, err := liteapi.NewClientWithDefaultTestnet()
	// tongoClient, err := liteapi.NewClient(liteapi.WithConfigurationFile(*config))
	fmt.Printf("main init\n")
	// 4IpHcYwLFsCq6YmEj21G+nas2d6W4szoYiqk2eYsY2k=
	// 0Q17OGptPHPyBhRwus9gYNx/K/2tiWY3RtOfHm2YnYc=
	// zJVFoCAcn6eLAZPGgyUirhQKnBj2sv0WZzAjlO9G4K0=
	tongoClient, err := liteapi.NewClient(liteapi.WithLiteServers([]config.LiteServer{{Host: "127.0.0.1:8803", Key: "0Q17OGptPHPyBhRwus9gYNx/K/2tiWY3RtOfHm2YnYc="}}))
	if err != nil {
		fmt.Printf("Unable to create tongo client: %v\n", err)
		return
	}
	fmt.Printf("liteapi.NewClient ok\n")
	// testTime := 10000 * time.Second

	// ticker := time.NewTicker(time.Second / time.Duration(targetQPS))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var wg sync.WaitGroup

	targetQPS := 5000
	for i := 0; i < 1; i++ {
		_ = <-ticker.C
		fmt.Printf("do a test. %d\n", i)

		resCh := make(chan ReqResult, targetQPS)

		// send request
		for i := 0; i < targetQPS; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resCh <- sendRequest(tongoClient)
			}()
		}

		// receive result
		wg.Add(1)
		go func() {
			defer wg.Done()
			totalElapsed := time.Duration(0)
			totalReqCount := 0
			totalFailedCount := 0
			for {
				select {
				case res := <-resCh:
					totalElapsed += res.elapsed
					totalReqCount += 1
					if res.err != nil {
						fmt.Printf("request error: %v\n", res.err)
						totalFailedCount += 1
					} else {
						fmt.Printf("success. elapsed: %v\n", res.elapsed)
					}
					if totalReqCount >= targetQPS {
						fmt.Printf("========= request count: %d. failed count: %d. failed rate: %.2f%%. avg cost: %vms.\n", totalReqCount, totalFailedCount, (float64(totalFailedCount)/float64(totalReqCount))*100, int(totalElapsed.Milliseconds())/totalReqCount)
						return
					}
				}
			}

		}()
	}

	wg.Wait()

}

var accountList = []ton.Address{
	tongo.MustParseAddress("0:A560DC2F757589A7152835573D23117910207DDEDEFA717DBC82D8CF932E31B2"),
	tongo.MustParseAddress("-1:CFB0B11350E4052FB60ED14CBF2FEDB265D7AC9EEB4C6BE0EEEE60AA2CE5DB3E"),
	tongo.MustParseAddress("-1:5555555555555555555555555555555555555555555555555555555555555555"),
	tongo.MustParseAddress("-1:0000000000000000000000000000000000000000000000000000000000000000"),
	tongo.MustParseAddress("0:A560DC2F757589A7152835573D23117910207DDEDEFA717DBC82D8CF932E31B2")}

func sendRequest(client *liteapi.Client) ReqResult {
	start := time.Now()
	_, err := client.GetAccountState(context.Background(), accountList[rand.Intn(len(accountList))].ID)
	elapsed := time.Since(start)

	return ReqResult{elapsed, err}
}
