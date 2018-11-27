package starter

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/updater/config"
	"github.com/skycoin/services/updater/pkg/checker/active"
	"github.com/skycoin/services/updater/pkg/checker/passive"
	"github.com/skycoin/services/updater/pkg/updater"
	loggerPkg "github.com/skycoin/services/updater/pkg/logger"
	"github.com/skycoin/services/updater/store/services"
	"github.com/pkg/errors"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("updater")
	ErrServiceNotFound = errors.New("service definition not found")
)

type Starter struct {
	activeCheckers  map[string]active.Fetcher
	passiveCheckers map[string]passive.Subscriber
	updaters        map[string]updater.Updater
	config config.Configuration
}

func New(conf config.Configuration) *Starter {
	s := &Starter{
		activeCheckers:  map[string]active.Fetcher{},
		passiveCheckers: map[string]passive.Subscriber{},
		updaters:        map[string]updater.Updater{},
		config: conf,
	}

	services.InitStorer("json")

	s.createUpdaters(conf)
	s.createCheckers(conf)

	return s
}

func (s *Starter) Start() {
	for _, checker := range s.activeCheckers {
		go checker.Start()
	}

	for _, checker := range s.passiveCheckers {
		go checker.Start()
	}
}

func (s *Starter) Update(service string) error {
	// get updater
	serviceConfig, ok := s.config.Services[service]
	if !ok {
		return ErrServiceNotFound
	}

	updater := s.updaters[serviceConfig.Updater]

	// Try update
		err := <- updater.Update(service, serviceConfig.CheckTag, loggerPkg.NewLogger(service))
		if err != nil {
			logger.Errorf("error on update %s", err)
		}
		return nil
}

func (s *Starter) Stop() {
	for _, checker := range s.activeCheckers {
		checker.Stop()
	}

	for _, checker := range s.passiveCheckers {
		checker.Stop()
	}
}

func (s *Starter) createUpdaters(conf config.Configuration) {
	for name, c := range conf.Updaters {
		u := updater.New(c.Kind, conf)
		s.updaters[name] = u
	}
}

func (s *Starter) createCheckers(conf config.Configuration) {
	for name, c := range conf.Services {
		if c.ActiveUpdateChecker != "" {
			activeConfig, ok := conf.ActiveUpdateCheckers[c.ActiveUpdateChecker]
			if !ok {
				logrus.Warnf("%s checker not defined for service %s, skipping service",
					c.ActiveUpdateChecker, name)
				continue
			}

			interval, err := time.ParseDuration(activeConfig.Interval)
			if err != nil {
				logrus.Fatalf("cannot parse interval %s of active checker configuration %s. %s", activeConfig.Interval,
					c.ActiveUpdateChecker, err)
			}
			log := loggerPkg.NewLogger(name)

			checker := active.New(name, &conf, s.updaters[c.Updater], log)
			checker.SetInterval(interval)
			s.activeCheckers[name] = checker
		} else {
			passiveConfig, ok := conf.PassiveUpdateCheckers[c.PassiveUpdateChecker]
			if !ok {
				logrus.Warnf("%s checker not defined for service %s, skipping service",
					c.ActiveUpdateChecker, name)
				continue
			}
			log := loggerPkg.NewLogger(name)

			sub := passive.New(passiveConfig, s.updaters[c.Updater], log)
			s.passiveCheckers[name] = sub
			sub.Subscribe(passiveConfig.Topic)
		}
	}
}
