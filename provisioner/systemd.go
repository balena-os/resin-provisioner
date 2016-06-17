package provisioner

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	pathLib "path"
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
	name := pathLib.Base(SUPERVISOR_PATH)
	return c.UnitStatus(name)
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

func (c *dbusConnection) EnableStartUnit(path string) error {
	name := pathLib.Base(path)
	ch := make(chan string)

	if _, _, err := c.EnableUnitFiles([]string{path}, false, false); err != nil {
		return err
	} else if _, err := c.StartUnit(name, "replace", ch); err != nil {
		return err
	}

	// Block until attempt to start unit succeeds/fails.
	if result := <-ch; result != "done" {
		return fmt.Errorf("Start failed due to %s.", result)
	}

	return nil
}

func (c *dbusConnection) RestartUnitNoWait(path string) error {
	name := pathLib.Base(path)
	_, err := c.RestartUnit(name, "replace", nil)

	return err
}

func (c *dbusConnection) SupervisorEnableStart() error {
	// We need to restart the prepare-openvpn.service to avoid a bug whereby
	// config.json is read before endpoints are populated, resulting in
	// misconfigured openvpn.
	if err := c.RestartUnitNoWait(PREPARE_OPENVPN_PATH); err != nil {
		return err
	} else if err := c.EnableStartUnit(SUPERVISOR_PATH); err != nil {
		return err
	}

	// We first need to ensure that the update-resin-supervisor.service is
	// able to perform the first supervisor image pull by setting the tag in
	// /etc/supervisor.conf.
	if err := setSupervisorTag(); err != nil {
		return err
	}

	// Start the resin update timer too.
	return c.EnableStartUnit(UPDATE_RESIN_TIMER_PATH)
}
