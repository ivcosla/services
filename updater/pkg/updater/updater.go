package updater

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/updater/config"
	"github.com/skycoin/services/updater/pkg/logger"
)

type Updater interface {
	Update(service, version string, log *logger.Logger) chan error
}

func New(kind string, conf config.Configuration) Updater {

	normalized := strings.ToLower(kind)
	logrus.Infof("updater: %s", normalized)

	switch normalized {
	case "custom":
		return newCustomUpdater(conf.Services)
	}

	return newCustomUpdater(conf.Services)
}
