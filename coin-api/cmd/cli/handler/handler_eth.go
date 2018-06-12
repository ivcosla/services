package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	minEthAddrLen = 26
	maxEthAddrLen = 35
)


func init() {
	client = &http.Client{
		Timeout: clientTimeout,
	}
}

// ETH is a cli bitcoin handler
type ETH struct{}

// NewETH returns new bitcoin handler instance
func NewETH() *ETH {
	return &ETH{}
}

// GenerateKeyPair generates keypair for bitcoin
func (e *ETH) GenerateKeyPair(c *cli.Context) error {
	req, err := http.NewRequest(http.MethodPost, "/keys", nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Key %s created\n", resp)
	return nil
}

// GenerateAddress generates addresses and keypairs for bitcoin
func (e *ETH) GenerateAddress(c *cli.Context) error {
	publicKey := c.Args().Get(1)

	params := map[string]interface{}{
		"publicKey": publicKey,
	}

	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, "/address", body)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Address %s created\n", resp)

	return nil
}

// CheckBalance checks bitcoin balance
func (e *ETH) CheckBalance(c *cli.Context) error {
	addr := c.Args().First()

	if len(addr) > 35 || len(addr) < 26 {
		err := errors.New(fmt.Sprintf("Address lenght must be between %d and %d",
			minEthAddrLen, maxEthAddrLen))
		return err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/address/%s", addr), nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Check balance success %s\n", resp)
	return nil
}
