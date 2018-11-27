package passive

import (
	"strings"

	"github.com/skycoin/services/updater/config"
	"github.com/skycoin/services/updater/pkg/logger"
	"github.com/skycoin/services/updater/pkg/updater"
)

type Subscriber interface {
	Subscribe(topic string)
	Start()
	Stop()
}

func New(config config.SubscriberConfig, updater updater.Updater, log *logger.Logger) Subscriber {
	config.MessageBroker = strings.ToLower(config.MessageBroker)
	switch config.MessageBroker {
	case "nats":
		return newNats(updater, config.Urls[0], config.NotifyUrl, log)
	}

	return newNats(updater, config.Urls[0], config.NotifyUrl, log)
}

