package provisioner

// TODO: Need to enable update-resin-supervisor.timer.

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
)

const (
	SUPERVISOR_NAME = "resin-supervisor.service"
	SUPERVISOR_PATH = "/lib/systemd/system/resin-supervisor.service"
)

type dbusConnection struct {
	*dbus.Conn
}

func NewDbus() (ret *dbusConnection, err error) {
	var conn *dbus.Conn

	if conn, err = dbus.New(); err == nil {
		ret = &dbusConnection{conn}
	}

	return
}

func (c *dbusConnection) UnitStatus(name string) (status *dbus.UnitStatus, err error) {
	var statuses []dbus.UnitStatus

	if statuses, err = c.ListUnitsByNames([]string{name}); err == nil {
		if len(statuses) == 0 {
			return
		}
		if len(statuses) > 1 {
			err = fmt.Errorf("%d units returned for name '%s', expected 1.",
				len(statuses), name)
		} else {
			status = &statuses[0]
		}
	}

	return
}

func (c *dbusConnection) SupervisorStatus() (*dbus.UnitStatus, error) {
	return c.UnitStatus(SUPERVISOR_NAME)
}

func (c *dbusConnection) SupervisorRunning() (ret bool, err error) {
	if status, err := c.SupervisorStatus(); err != nil {
		return false, err
	} else {
		running := status != nil &&
			status.LoadState == "loaded" &&
			status.ActiveState == "active"

		return running, nil
	}
}
