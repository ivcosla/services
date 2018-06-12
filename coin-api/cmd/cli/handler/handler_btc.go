package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/skycoin/services/coin-api/internal/server"
	"github.com/skycoin/services/coin-api/internal/btc"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	apiVersion = "api/v1"
	serverAddress = ":9090"
	btcEndpoint = "btc"
	clientTimeout = time.Second * 10
	minBtcAddrLen = 26
	maxBtcAddrLen = 35
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{
		Timeout: clientTimeout,
	}
}

// BTC is a cli bitcoin handler
type BTC struct{}

// NewBTC returns new bitcoin handler instance
func NewBTC() *BTC {
	return &BTC{}
}

// GenerateKeyPair generates keypair for bitcoin
func (b *BTC) GenerateKeyPair(c *cli.Context) error {
	req, err := http.NewRequest(http.MethodPost, serverAddress + apiVersion + btcEndpoint + "/keys", nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	keyPairResponse := &server.KeyPairResponse{}

	json.NewDecoder(resp.Body).Decode(keyPairResponse)
	log.Printf("Key %s created\n", keyPairResponse)
	return nil
}

// GenerateAddress generates addresses and keypairs for bitcoin
func (b *BTC) GenerateAddress(c *cli.Context) error {
	publicKey := c.Args().Get(1)

	params := map[string]interface{}{
		"key": publicKey,
	}

	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, serverAddress + apiVersion + btcEndpoint + "/address", body)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	
	addressResponse := &server.AddressResponse{}
	json.NewDecoder(resp.Body).Decode(addressResponse)

	log.Printf("Address %s created\n", addressResponse)

	return nil
}

// CheckBalance checks bitcoin balance
func (b *BTC) CheckBalance(c *cli.Context) error {
	addr := c.Args().First()

	if len(addr) > maxBtcAddrLen || len(addr) < minBtcAddrLen {
		err := errors.New(fmt.Sprintf("Address length must be between %d and %d",
			minBtcAddrLen, maxBtcAddrLen))
		return err
	}
 
	req, err := http.NewRequest(http.MethodGet, serverAddress + apiVersion + btcEndpoint + fmt.Sprintf("/address/%s", addr), nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	balanceResponse := &btc.BalanceResponse{}
	json.NewDecoder(resp.Body).Decode(balanceResponse)

	log.Printf("Check balance success %+v\n", balanceResponse)
	return nil
}

func (b *BTC) CheckTransaction(c *cli.Context) error {
	txId := c.Args().First()

	req, err := http.NewRequest(http.MethodGet, serverAddress + apiVersion + btcEndpoint + fmt.Sprintf("/transaction/%s", txId), nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	txStatusResponse := &btc.TxStatus{}
	json.NewDecoder(resp.Body).Decode(txStatusResponse)

	log.Printf("Check balance success %+v\n", txStatusResponse)
	return nil
}
