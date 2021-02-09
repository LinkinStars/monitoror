package notification

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/notification/dingding"
	"github.com/monitoror/monitoror/store"
)

const (
	MaxAlertTimes      = 3
	QueryName          = "notification-name"
	notificationFormat = `# [monitoror warning]
**%s** exception occur
> details: %s %s %s`
)

var (
	alertRecord = make(map[string]int, 0)
	mutex       sync.Mutex
	notifiers   = make(map[coreModels.NotificationType]coreModels.Notification, 0)
)

// RegisterNotification register all notification handler
func RegisterNotification(s *store.Store) {
	notifiers[coreModels.DingDingNotice] = dingding.NewDingDingRobotNotifier(s)

	for notifierType, notifier := range notifiers {
		for _, name := range notifier.GetVariantsNames() {
			if validate, _ := notifier.Validate(name); !validate {
				delete(notifiers, notifierType)
			} else {
				notifier.Enable(name)
			}
		}
	}
}

// ResponseInterceptor intercept the response and process the notification
func ResponseInterceptor() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		noticeName := c.QueryParam(QueryName)
		tile := &coreModels.Tile{}
		err := json.Unmarshal(resBody, tile)
		if err != nil {
			return
		}

		go func() {
			dealNotification(tile, noticeName)
		}()
	})
}

func dealNotification(tile *coreModels.Tile, noticeName string) {
	mutex.Lock()
	defer mutex.Unlock()

	// TODO Not all situation is success, need to judge based on the request
	if tile.Status == coreModels.SuccessStatus {
		alertRecord[tile.Label] = 0
	} else if alertRecord[tile.Label] <= MaxAlertTimes {
		alertRecord[tile.Label]++
	}

	// when the maximum number of notifications has been exceeded, this notification will be discarded
	if 0 < alertRecord[tile.Label] && alertRecord[tile.Label] <= MaxAlertTimes {
		for _, notifier := range notifiers {
			_ = notifier.Notify(fmt.Sprintf(notificationFormat, noticeName, tile.Label, tile.Message, tile.Status))
		}
	}
}
