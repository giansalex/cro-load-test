module github.com/giansalex/cro-load-test

go 1.15

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

require (
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/gorilla/websocket v1.4.2
	github.com/prometheus/client_golang v1.9.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.3
	github.com/tendermint/tendermint v0.34.4
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)
