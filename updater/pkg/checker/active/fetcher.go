package active

import (
	"time"

	"github.com/skycoin/services/updater/config"
	"github.com/skycoin/services/updater/pkg/logger"
	"github.com/skycoin/services/updater/pkg/updater"
)

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}


func New(service string, c *config.Configuration, updater updater.Updater, log *logger.Logger) Fetcher {
	serviceConfig := c.Services[service]
	updateCheckerConfig := c.ActiveUpdateCheckers[serviceConfig.ActiveUpdateChecker]
	updaterConfig := c.Updaters[serviceConfig.Updater]

	if updaterConfig.Kind == "swarm" {
		log.Info("Swarm mode cannot fetch from Git, falling back to Dockerhub")
		updateCheckerConfig.Kind = "dockerhub"
	}

	switch updateCheckerConfig.Kind {
	case "git":
		return newGit(updater, service, serviceConfig.Repository, updateCheckerConfig.NotifyUrl, log)
	case "naive":
		return newNaive(updater, service, serviceConfig.Repository, updateCheckerConfig.NotifyUrl, log)
	}
	return newNaive(updater, service, serviceConfig.Repository, updateCheckerConfig.NotifyUrl, log)
}
