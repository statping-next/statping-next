package notifications

import (
	"github.com/statping-next/statping-next/types/metrics"
)

func (n *Notification) AfterFind() (err error) {
	metrics.Query("notifier", "find")
	return
}

func (n *Notification) AfterCreate() {
	metrics.Query("notifier", "create")
}

func (n *Notification) AfterUpdate() {
	metrics.Query("notifier", "update")
}

func (n *Notification) AfterDelete() {
	metrics.Query("notifier", "delete")
}
