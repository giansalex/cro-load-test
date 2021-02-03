package loadtest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

const (
	accounName = "golang"
)

type Signature struct {
	keyBase            keyring.Keyring
	interfacesRegistry codectypes.InterfaceRegistry
}

// RegisterInterfaces register decoding interface to the decoder by using the provided interface
// registry.
func (signature *Signature) RegisterInterfaces(registry func(registry codectypes.InterfaceRegistry)) *Signature {
	registry(signature.interfacesRegistry)

	return signature
}

func (signature *Signature) Import(armor, pass string) error {
	kb := keyring.NewInMemory()

	err := kb.ImportPrivKey("golang", armor, pass)
	if err != nil {
		return err
	}

	signature.keyBase = kb

	return nil
}

func (signature *Signature) Recover(paraphrase string) (keyring.Info, error) {
	kb := keyring.NewInMemory()

	hdPath := hd.CreateHDPath(cosmostypes.GetConfig().GetCoinType(), 0, 0)

	info, err := kb.NewAccount(accounName, paraphrase, "", hdPath.String(), hd.Secp256k1)
	if err != nil {
		return nil, err
	}

	signature.keyBase = kb

	return info, nil
}

func (signature *Signature) GetTxConfig() client.TxConfig {
	marshaler := codec.NewProtoCodec(signature.interfacesRegistry)

	return authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)
}

func (signature *Signature) Sign(accNro, sequence uint64, txBuilder client.TxBuilder) ([]byte, error) {

	marshaler := codec.NewProtoCodec(signature.interfacesRegistry)
	txConfig := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	signerData := authsigning.SignerData{
		ChainID:       "crossfire",
		AccountNumber: accNro,
		Sequence:      sequence,
	}

	key, err := signature.keyBase.Key(accounName)
	if err != nil {
		return nil, err
	}
	pubKey := key.GetPubKey()

	signMode := signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	sigBytes, _, err := signature.keyBase.Sign(accounName, bytesToSign)
	if err != nil {
		return nil, err
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: sequence,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	parsedTx := txBuilder.GetTx()

	data, err := authtx.DefaultTxEncoder()(parsedTx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (signature *Signature) ParseJson(json string) (cosmostypes.Tx, error) {

	marshaler := codec.NewProtoCodec(signature.interfacesRegistry)

	return authtx.DefaultJSONTxDecoder(marshaler)([]byte(json))
}

// NewDecoder creates a new decoder
func NewSignature() *Signature {
	interfaceRegistry := codectypes.NewInterfaceRegistry()

	return &Signature{
		interfacesRegistry: interfaceRegistry,
	}
}

// DefaultDecoder is a decoder with all Cosmos builtin modules interfaces registered
var DefaultSignature = NewSignature().RegisterInterfaces(RegisterDefaultInterfaces)
