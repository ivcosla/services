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
	ETH_ADDRESS_LENGTH = 40
	
)


func init() {
	client = &http.Client{
		Timeout: clientTimeout,
	}
}

// ETH is a cli ethereum handler
type ETH struct{}

// NewETH returns new ethereum handler instance
func NewETH() *ETH {
	return &ETH{}
}

// GenerateKeyPair generates keypair for ethereum
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

// GenerateAddress generates addresses and keypairs for ethereum
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

// CheckBalance checks ethereum balance
func (e *ETH) CheckBalance(c *cli.Context) error {
	addr := c.Args().First()

	if len(addr) != ETH_ADDRESS_LENGTH { 
		err := errors.New(fmt.Sprintf("Address lenght must be %d",
			ETH_ADDRESS_LENGTH))
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
