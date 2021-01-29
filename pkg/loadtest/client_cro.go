package loadtest

import (
	"errors"
	"fmt"
)

// MyABCIAppClientFactory creates instances of MyABCIAppClient
type MyABCIAppClientFactory struct {
	account    uint64
	sequence   uint64
	paraphrase string
	max        uint64
}

// MyABCIAppClientFactory implements loadtest.ClientFactory
var _ ClientFactory = (*MyABCIAppClientFactory)(nil)

const jsonMsg = "{\"body\":{\"messages\":[{\"@type\":\"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress\",\"delegator_address\":\"cro1lglgwsrt7m293me6kgwh4vnw55yrgulggss8t5\",\"withdraw_address\":\"cro1lglgwsrt7m293me6kgwh4vnw55yrgulggss8t5\"}],\"memo\":\"\",\"timeout_height\":\"0\",\"extension_options\":[],\"non_critical_extension_options\":[]},\"auth_info\":{\"signer_infos\":[],\"fee\":{\"amount\":[{\"denom\":\"basetcro\",\"amount\":\"20000\"}],\"gas_limit\":\"200000\",\"payer\":\"\",\"granter\":\"\"}},\"signatures\":[]}"

// MyABCIAppClient is responsible for generating transactions. Only one client
// will be created per connection to the remote Tendermint RPC endpoint, and
// each client will be responsible for maintaining its own state in a
// thread-safe manner.
type MyABCIAppClient struct {
	txs        map[uint64][]byte
	sequence   uint64
	count      uint64
	max        uint64
	paraphrase string
}

// MyABCIAppClient implements loadtest.Client
var _ Client = (*MyABCIAppClient)(nil)

func NewABCIAppClientFactory(paraphrase string) *MyABCIAppClientFactory {
	return &MyABCIAppClientFactory{paraphrase: paraphrase}
}

func (f *MyABCIAppClientFactory) ValidateConfig(cfg Config) error {
	// Do any checks here that you need to ensure that the load test
	// configuration is compatible with your client.
	return nil
}

func (f *MyABCIAppClientFactory) NewClient(cfg Config) (Client, error) {
	signer := DefaultSignature
	_, err := signer.Recover(f.paraphrase)
	if err != nil {
		return nil, err
	}
	totalTxs := f.max + 2

	fmt.Printf("Total Txs signed: %d \n", totalTxs)
	if totalTxs < 1 {
		return nil, errors.New("Invalid total txs")
	}
	msgTx, err := signer.ParseJson(jsonMsg)
	if err != nil {
		return nil, err
	}

	txs := make(map[uint64][]byte, totalTxs)
	var i uint64
	for i = 0; i < totalTxs; i++ {
		seq := f.sequence + uint64(i)
		data, err := signer.Sign(f.account, seq, msgTx)
		if err != nil {
			return nil, err
		}

		txs[seq] = data
		// fmt.Println("sequence sign:", seq)
	}

	return &MyABCIAppClient{txs: txs, sequence: f.sequence, count: 0, max: f.max}, nil
}

// GenerateTx must return the raw bytes that make up the transaction for your
// ABCI app. The conversion to base64 will automatically be handled by the
// loadtest package, so don't worry about that. Only return an error here if you
// want to completely fail the entire load test operation.
func (c *MyABCIAppClient) GenerateTx() ([]byte, error) {
	if c.count >= c.max {
		return nil, errors.New("---Max Tx limit---")
	}
	seq := c.sequence + c.count
	c.count++
	// fmt.Println("sequence send:", seq)
	return c.txs[seq], nil
}
