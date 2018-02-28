package model

import (
	"container/list"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin/services/otc/types"
)

var (
	ErrUnknownStatus = errors.New("unknown status type")
	ErrNilService    = errors.New("nil service passed to model")
)

type Model struct {
	sync.Mutex

	path    string
	stop    chan struct{}
	storage *Storage
	lookup  *Lookup
	results *list.List
	config  *types.Config
	logger  *log.Logger
	errs    *log.Logger
	events  *os.File
	Scanner types.Service
	Sender  types.Service
	Monitor types.Service
}

func NewModel(c *types.Config, scn, sndr, mntr types.Service, errs *log.Logger) (*Model, error) {
	m := &Model{
		lookup:  NewLookup(),
		results: list.New().Init(),
		path:    c.Model.Path,
		stop:    make(chan struct{}),
		config:  c,
		logger:  log.New(os.Stdout, types.LOG_MODEL, types.LOG_FLAGS),
		errs:    errs,
		Scanner: scn,
		Sender:  sndr,
		Monitor: mntr,
	}

	// make sure all services are there
	if scn == nil || sndr == nil || mntr == nil {
		return nil, ErrNilService
	}

	// make sure db dir exists
	_, err := os.Stat(m.path)
	if err != nil {
		return nil, err
	}

	// open request storage struct
	if m.storage, err = NewStorage(c.Model.Path); err != nil {
		return nil, err
	}

	// open events log file
	if m.events, err = os.OpenFile(
		m.path+"events.json",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	); err != nil {
		return nil, err
	}

	// get list of files in db dir
	files, err := ioutil.ReadDir(m.path + "requests/")
	if err != nil {
		return nil, err
	}

	// for each .json file in db dir
	for _, file := range files {
		// create a slice of requests contained in file
		requests, err := m.storage.LoadRequests(file.Name())
		if err != nil {
			if err == io.EOF {
				continue
			}
			return nil, err
		}

		// inject each request into the proper service
		for _, request := range requests {
			if err := m.Add(request); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m *Model) Stop() {
	m.Scanner.Stop()
	m.Sender.Stop()
	m.Monitor.Stop()

	m.stop <- struct{}{}
	m.logger.Println("stopped")
}

func (m *Model) Start() {
	m.logger.Println("started")
	go func() {
		for {
			<-time.After(time.Second * time.Duration(m.config.Model.Tick))

			select {
			case <-m.stop:
				return
			default:
				m.process()
			}
		}
	}()
}

func (m *Model) process() {
	m.Lock()
	defer m.Unlock()

	for e := m.results.Front(); e != nil; e = e.Next() {
		// convert to result promise
		r := e.Value.(chan *types.Result)

		// non-blocking read on each result promise
		select {
		case result := <-r:
			if result.Err != nil {
				// TODO: re-route request, try again?
				m.errs.Printf("model: %v\n", result.Err)
			} else {
				result.Request.Metadata.Update()

				// save new state to disk
				if err := m.storage.SaveRequest(result.Request); err != nil {
					m.errs.Printf("model: %v\n", result.Err)
				}

				// send to next service if request isn't finished
				if next := m.Handle(result.Request); next != nil {
					// add result promise to queue
					m.results.PushBack(next)
				}
			}

			// append to events log
			if err := NewEvent(
				result.Request,
				result.Err,
			).Append(m.events); err != nil {
				m.errs.Printf("model: %v\n", err)
			}

			// this elem has been handled, so remove
			m.results.Remove(e)
		default:
			continue
		}
	}
}

func (m *Model) Add(request *types.Request) error {
	m.Lock()
	defer m.Unlock()

	// save to disk
	if err := m.storage.SaveRequest(request); err != nil {
		return err
	}

	m.lookup.SetDrop(request.Drop, request.Currency, request.Address)

	// route to next component
	if result := m.Handle(request); result != nil {
		// add to end of queue
		m.results.PushBack(result)
	}

	return nil
}

func (m *Model) Handle(r *types.Request) chan *types.Result {
	switch r.Metadata.Status {
	case types.DEPOSIT:
		return m.Scanner.Handle(r)
	case types.SEND:
		return m.Sender.Handle(r)
	case types.CONFIRM:
		return m.Monitor.Handle(r)
	case types.EXPIRED:
		fallthrough
	case types.DONE:
		fallthrough
	default:
		return nil
	}
}

func (m *Model) GetMetadata(drop types.Drop, curr types.Currency) (*types.Metadata, error) {
	// lookup sky address for filename
	address, err := m.lookup.GetAddress(drop, curr)
	if err != nil {
		return nil, err
	}

	// get metadata from disk
	metadata, err := m.storage.LoadMetadata(address, drop, curr)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
