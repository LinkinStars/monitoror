package dingding

import (
	"github.com/LinkinStars/dingrobot"

	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	dingdingConfig "github.com/monitoror/monitoror/notification/dingding/config"
	"github.com/monitoror/monitoror/store"
)

// Notification ding ding robot notifier
type Notification struct {
	store *store.Store

	config map[coreModels.VariantName]*dingdingConfig.DingDing

	robot   dingrobot.Roboter
	message chan string
}

// NewDingDingRobotNotifier create a dingding robot notifier
func NewDingDingRobotNotifier(store *store.Store) (n *Notification) {
	n = &Notification{}
	n.store = store
	n.config = make(map[coreModels.VariantName]*dingdingConfig.DingDing)
	n.message = make(chan string)

	// Load core config from env
	pkgMonitorable.LoadConfig(&n.config, dingdingConfig.Default)

	return n
}

func (n *Notification) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(n.config)
}

func (n *Notification) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := n.config[variantName]

	// No configuration set
	if conf.Webhook == "" && conf.Secret == "" {
		return false, nil
	}

	// Validate Config
	if errors := pkgMonitorable.ValidateConfig(conf, variantName); errors != nil {
		return false, errors
	}

	return true, nil
}

func (n *Notification) Enable(variantName coreModels.VariantName) {
	conf := n.config[variantName]
	if len(conf.Webhook) == 0 {
		return
	}

	// set basic config
	n.robot = dingrobot.NewRobot(conf.Webhook)
	if len(conf.Secret) > 0 {
		n.robot.SetSecret(conf.Secret)
	}

	n.Handle()
	return
}

// Notify send ding ding message
func (n *Notification) Notify(content string) (err error) {
	n.message <- content
	return
}

// Handle start handle notification
func (n *Notification) Handle() {
	go func() {
		for m := range n.message {
			// TODO deal error and set title more elegant, maybe we need to add config for @ someone or @ all
			_ = n.robot.SendMarkdown("[monitoror warning]", m, []string{}, false)
		}
	}()
}
