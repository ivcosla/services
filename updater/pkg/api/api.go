package api

import (
	"net/http"
	"strings"
	"context"
	"net"
	"net/url"
	"github.com/skycoin/services/updater/pkg/starter"
	"github.com/skycoin/services/updater/config"
	"encoding/json"
	"github.com/skycoin/skycoin/src/util/logging"
)

var logger = logging.MustGetLogger("api")

// Gateway represents an API to communicate with the node
type Gateway interface {
	Start(string) error
	Stop() error
}

// HTTPResponse represents the http response struct
type HTTPResponse struct {
	Error *HTTPError  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// HTTPError is included in an HTTPResponse
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewHTTPErrorResponse returns an HTTPResponse with the Error field populated
func NewHTTPErrorResponse(code int, msg string) HTTPResponse {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return HTTPResponse{
		Error: &HTTPError{
			Code:    code,
			Message: msg,
		},
	}
}

// ClientError is used for non-200 API responses
type ClientError struct {
	Status     string
	StatusCode int
	Message    string
}

// NewClientError creates a ClientError
func NewClientError(status string, statusCode int, message string) ClientError {
	return ClientError{
		Status:     status,
		StatusCode: statusCode,
		Message:    strings.TrimRight(message, "\n"),
	}
}

func (e ClientError) Error() string {
	return e.Message
}

// ServerGateway implements gateway interface for REST server
type ServerGateway struct {
	server *http.Server
	starter *starter.Starter
}

// NewServerGateway returns a ServerGateway
func NewServerGateway(conf config.Configuration) *ServerGateway {
	return &ServerGateway{
		starter: starter.New(conf),
	}
}

// Start starts the REST server gateway
func (s *ServerGateway) Start(addrs string) error {
	l, err := net.Listen("tcp", addrs)
	if err != nil {
		return err
	}

	s.server = &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", s.Update)
	s.server.Handler = mux

	s.starter.Start()
	return s.server.Serve(l)
}

// Stop closes the REST server gateway
func (s *ServerGateway) Stop() error {
	s.starter.Stop()
	return s.server.Shutdown(context.Background())
}

// Update gets the service that needs to be updated and updates it
// URI: /update/:service_name
func (s *ServerGateway) Update(w http.ResponseWriter, r *http.Request) {
	service := retrieveServiceFromURL(r.URL)
	err := s.starter.Update(service)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError,
			NewHTTPErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	} else {
		writeJSON(w, http.StatusOK, HTTPResponse{Data:service + " updated"})
	}

}

// retrievePkFromURL returns the id used on endpoints of the form path/:pk
// it doesn't checks if the endpoint has this form and can fail with other
// endpoint forms
func retrieveServiceFromURL(url *url.URL) string {
	splittedPath := strings.Split(url.EscapedPath(), "/")
	return splittedPath[len(splittedPath)-1]
}

// writeJSON writes a json object on a http.ResponseWriter with the given code,
// panics on marshaling error
func writeJSON(w http.ResponseWriter, code int, object interface{}) {
	jsonObject, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonObject)
	if err != nil {
		logger.Error(err)
	}
}
