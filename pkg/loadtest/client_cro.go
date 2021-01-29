package loadtest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
)

// MyABCIAppClientFactory creates instances of MyABCIAppClient
type MyABCIAppClientFactory struct {
	paraphrase string
}

// MyABCIAppClientFactory implements loadtest.ClientFactory
var _ ClientFactory = (*MyABCIAppClientFactory)(nil)

const jsonMsg = "{\"body\":{\"messages\":[{\"@type\":\"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress\",\"delegator_address\":\"cro1lglgwsrt7m293me6kgwh4vnw55yrgulggss8t5\",\"withdraw_address\":\"cro1lglgwsrt7m293me6kgwh4vnw55yrgulggss8t5\"}],\"memo\":\"\",\"timeout_height\":\"0\",\"extension_options\":[],\"non_critical_extension_options\":[]},\"auth_info\":{\"signer_infos\":[],\"fee\":{\"amount\":[{\"denom\":\"basetcro\",\"amount\":\"20000\"}],\"gas_limit\":\"200000\",\"payer\":\"\",\"granter\":\"\"}},\"signatures\":[]}"

// MyABCIAppClient is responsible for generating transactions. Only one client
// will be created per connection to the remote Tendermint RPC endpoint, and
// each client will be responsible for maintaining its own state in a
// thread-safe manner.
type MyABCIAppClient struct {
	txs     map[uint64][]byte
	lcd     *LcdClient
	signer  *Signature
	keyInfo keyring.Info
	msg     cosmostypes.Tx
	count   uint64
	max     uint64
	seq     uint64
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
	info, err := signer.Recover(f.paraphrase)
	if err != nil {
		return nil, err
	}
	fmt.Println("Wallet Address:", info.GetAddress().String())
	msgTx, err := signer.ParseJson(jsonMsg)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	lcd := NewLcdClient(client, cfg.LcdEndpoint)

	return &MyABCIAppClient{signer: signer, keyInfo: info, msg: msgTx, lcd: lcd, max: 10, count: 0}, nil
}

// GenerateTx must return the raw bytes that make up the transaction for your
// ABCI app. The conversion to base64 will automatically be handled by the
// loadtest package, so don't worry about that. Only return an error here if you
// want to completely fail the entire load test operation.
func (c *MyABCIAppClient) GenerateTx() ([]byte, error) {

	if c.count >= c.max {
		c.count = 0
	}

	if c.count == 0 {
		err := c.makeTxs()
		if err != nil {
			return nil, err
		}
	}
	seq := c.seq + c.count
	c.count++

	return c.txs[seq], nil
}

func (c *MyABCIAppClient) makeTxs() error {
	account, err := c.lcd.Account(c.keyInfo.GetAddress().String())
	if err != nil {
		return err
	}

	totalTxs := c.max
	accountNro, _ := strconv.ParseUint(account.Result.Value.AccountNumber, 10, 64)
	sequence, _ := strconv.ParseUint(account.Result.Value.Sequence, 10, 64)

	txs := make(map[uint64][]byte, totalTxs)
	var i uint64
	for i = 0; i < c.max; i++ {
		seq := sequence + i
		data, err := c.signer.Sign(accountNro, seq, c.msg)
		if err != nil {
			return err
		}

		txs[seq] = data
	}

	c.txs = txs
	c.seq = sequence

	return nil
}
