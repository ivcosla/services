package admin

import (
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func New(curs *currencies.Currencies, modl *model.Model) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", Status(curs, modl))
	mux.HandleFunc("/api/pause", Pause(curs, modl))
	mux.HandleFunc("/api/price", Price(curs, modl))
	mux.HandleFunc("/api/source", Source(curs, modl))
	mux.HandleFunc("/api/transactions", Transactions(curs, modl))
	mux.HandleFunc("/api/transactions/pending", TransactionsPending(curs, modl))
	mux.HandleFunc("/api/transactions/completed", TransactionsCompleted(curs, modl))
	mux.HandleFunc("/api/addresses/sky", Addresses(otc.SKY, curs, modl))
	mux.HandleFunc("/api/holding/btc", Holding(otc.BTC, curs, modl))
	return mux
}
