package server

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/skycoin/skycoin/src/cipher"

	"github.com/pkg/errors"

	"github.com/skycoin/services/coin-api/internal/btc"
)

type KeyPairResponse struct {
	Public  string `json:"public"`
	Private []byte `json:"private"`
}

type AddressRequest struct {
	PublicKey string `json:"key"`
}

type AddressResponse struct {
	Address string `json:"address"`
}

type handlerBTC struct {
	btcService *btc.ServiceBtc
	checker    Checker
}

type BtcStats struct {
	NodeStatus string `json:"node_status"`
	NodeHost   string `json:"node_host"`
}

func newHandlerBTC(blockExplorer string, watcherUrl string) (*handlerBTC, error) {
	log.Printf("Start new BTC handler with watcher %s explorer %s", watcherUrl, blockExplorer)
	service, err := btc.NewBTCService(blockExplorer, watcherUrl)

	if err != nil {
		return nil, err
	}

	return &handlerBTC{
		btcService: service,
		checker:    service,
	}, nil
}

func (h *handlerBTC) generateKeyPair(ctx echo.Context) error {
	buffer := make([]byte, 256)
	_, err := rand.Read(buffer)

	if err != nil {
		return err
	}

	public, private := btc.ServiceBtc{}.GenerateKeyPair()

	if err := public.Verify(); err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result KeyPairResponse `json:"result"`
	}{
		"Ok",
		http.StatusOK,
		KeyPairResponse{
			Public:  public.Hex(),
			Private: private[:],
		},
	}

	// Write response with newly created key pair
	ctx.JSONPretty(http.StatusCreated, resp, "\t")
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	var req AddressRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(req.PublicKey)

	if err != nil {
		return handleError(ctx, err)
	}

	address, err := btc.ServiceBtc{}.GenerateAddr(publicKey)

	if err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result AddressResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: AddressResponse{
			Address: address,
		},
	}

	ctx.JSONPretty(http.StatusCreated, resp, "\t")
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	txId := ctx.Param("transid")

	resultChan := make(chan *btc.TxStatus)
	errChan := make(chan error)

	go func() {
		result, err := h.checker.CheckTxStatus(txId)
		status, ok := result.(*btc.TxStatus)

		if err != nil {
			errChan <- err
			return
		}

		if !ok {
			errChan <- errors.New("cannot convert result to *TxStatus")
			return
		}

		resultChan <- status
	}()

	var (
		status *btc.TxStatus
		err    error
		done   bool
	)
	select {
	case status = <-resultChan:
	case err = <-errChan:
	case <-ctx.Request().Context().Done():
		done = true
		log.Println("Request is canceled")
	}

	if done {
		return ctx.NoContent(http.StatusNoContent)
	}

	if err != nil {
		return handleError(ctx, err)
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string        `json:"status"`
		Code   int           `json:"code"`
		Result *btc.TxStatus `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: status,
	}, "\t")

	return nil
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	address := ctx.Param("address")

	resultChan := make(chan *btc.BalanceResponse)
	errChan := make(chan error)

	go func() {
		result, err := h.checker.CheckBalance(address)

		if err != nil {
			errChan <- err
			return
		}

		var (
			balance *btc.BalanceResponse
			ok      bool
		)

		balance, ok = result.(*btc.BalanceResponse)

		if !ok {
			errChan <- errors.New("cannot convert result to type *btc.BalanceResponse")
			return
		}

		resultChan <- balance
	}()

	var (
		balance *btc.BalanceResponse
		err     error
		done    bool
	)

	select {
	case balance = <-resultChan:
	case err = <-errChan:
	case <-ctx.Request().Context().Done():
		done = true
		log.Println("Request is canceled")
	}

	if done {
		return ctx.NoContent(http.StatusNoContent)
	}

	if err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string               `json:"status"`
		Code   int                  `json:"code"`
		Result *btc.BalanceResponse `json:"result"`
	}{
		Status: "Ok",
		Code:   http.StatusOK,
		Result: balance,
	}

	ctx.JSONPretty(http.StatusOK, resp, "\t")
	return nil
}

// Hook for collecting stats
func (h handlerBTC) CollectStatuses(stats *Status) {
	stats.Lock()
	defer stats.Unlock()
	stats.Stats["btc"] = &BtcStats{
		NodeStatus: h.btcService.GetStatus(),
		NodeHost:   h.btcService.WatcherHost(),
	}
}

func handleError(ctx echo.Context, err error) error {
	return ctx.JSONPretty(http.StatusOK, struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: err.Error(),
	}, "\t")
}
