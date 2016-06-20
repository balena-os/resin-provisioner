package provisioner

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	pathLib "path"

	"github.com/resin-os/resin-provisioner/util"
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

func (c *dbusConnection) EnableUnitsAsync(paths []string) error {
	_, _, err := c.EnableUnitFiles(paths, false, false)

	return err
}

func (c *dbusConnection) EnableStartUnit(path string) error {
	name := pathLib.Base(path)
	ch := make(chan string)

	if err := c.EnableUnitsAsync([]string{path}); err != nil {
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

func (c *dbusConnection) EnableResinServices() error {
	if services, err := util.ReadLines(RESIN_SERVICES_PATH); err != nil {
		return err
	} else {
		paths := make([]string, len(services))

		for i, service := range services {
			paths[i] = fmt.Sprintf("%s%s", SERVICES_ROOT_PATH, service)
		}

		return c.EnableUnitsAsync(paths)
	}
}

func (c *dbusConnection) SupervisorEnableStart() error {
	// Start by enabling all required services.
	if err := c.EnableResinServices(); err != nil {
		return err
	}

	// We need to ensure that the update-resin-supervisor.service is able to
	// perform the first supervisor image pull by setting the tag in
	// /etc/supervisor.conf.
	if err := setSupervisorTag(); err != nil {
		return err
	}
	// Next, trigger the first update-resin-supervisor to make sure the
	// supervisor image is downloaded.
	if err := c.EnableStartUnit(UPDATE_RESIN_PATH); err != nil {
		return err
	}
	// We need to restart the prepare-openvpn.service ('wanted' by
	// openvpn-resin.service) to avoid a bug whereby config.json is read
	// before endpoints are populated, resulting in misconfigured openvpn.
	if err := c.RestartUnitNoWait(OPENVPN_PATH); err != nil {
		return err
	}

	// Start the resin update timer too.
	return c.EnableStartUnit(UPDATE_RESIN_TIMER_PATH)
}
