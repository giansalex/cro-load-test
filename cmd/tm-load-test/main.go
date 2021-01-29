package main

import (
	"log"
	"os"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/giansalex/cro-load-test/pkg/loadtest"
)

const appLongDesc = `Load testing application for Tendermint with optional master/slave mode.
Generates large quantities of arbitrary transactions and submits those 
transactions to one or more Tendermint endpoints. By default, it assumes that
you are running the Crossfire Crypto.com chain on your Tendermint network.

To run the application in a similar fashion to cro-bench (STANDALONE mode):
    cro-load-test -c 1 -T 10 -r 1000 -s 250 \
        --broadcast-tx-method async \
        --endpoints ws://tm-endpoint1.somewhere.com:26657/websocket,ws://tm-endpoint2.somewhere.com:26657/websocket

To run the application in MASTER mode:
    cro-load-test \
        master \
        --expect-slaves 2 \
        --bind localhost:26670 \
        --shutdown-wait 60 \
        -c 1 -T 10 -r 1000 -s 250 \
        --broadcast-tx-method async \
        --endpoints ws://tm-endpoint1.somewhere.com:26657/websocket,ws://tm-endpoint2.somewhere.com:26657/websocket

To run the application in SLAVE mode:
    cro-load-test slave --master localhost:26680

NOTES:
* MASTER mode exposes a "/metrics" endpoint in Prometheus plain text format
  which shows total number of transactions and the status for the master and
  all connected slaves.
* The "--shutdown-wait" flag in MASTER mode is specifically to allow your 
  monitoring system some time to obtain the final Prometheus metrics from the
  metrics endpoint.
* In SLAVE mode, all load testing-related flags are ignored. The slave always 
  takes instructions from the master node it's connected to.
`

func main() {

	wallet := os.Getenv("MAX")
	if wallet == "" {
		log.Fatal("Required WALLET environment")
	}

	configCro()

	appFactory := loadtest.NewABCIAppClientFactory(wallet)

	if err := loadtest.RegisterClientFactory("cro-crossfire", appFactory); err != nil {
		panic(err)
	}

	loadtest.Run(&loadtest.CLIConfig{
		AppName:              "cro-load-test",
		AppShortDesc:         "Load testing application for Crypto.com Chain",
		AppLongDesc:          appLongDesc,
		DefaultClientFactory: "cro-crossfire",
	})
}

func configCro() {
	config := cosmostypes.GetConfig()
	config.SetBech32PrefixForAccount("cro", "cropub")
	config.SetBech32PrefixForValidator("crocncl", "crocnclpub")
	config.SetBech32PrefixForConsensusNode("crocnclcons", "crocnclconspub")
	config.SetCoinType(394)                         // required by sign
	config.SetFullFundraiserPath("44'/394'/0'/0/1") // required by sign

	config.Seal()
}
