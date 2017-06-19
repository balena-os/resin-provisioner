package provisioner

import (
	"fmt"
	pathLib "path"

	"github.com/coreos/go-systemd/dbus"
)

type dbusConnection struct {
	*dbus.Conn
}

func newDbus() (ret *dbusConnection, err error) {
	var conn *dbus.Conn

	if conn, err = dbus.New(); err == nil {
		ret = &dbusConnection{conn}
	}

	return
}

func (c *dbusConnection) restartUnit(path string) error {
	name := pathLib.Base(path)

	if _, err := c.RestartUnit(name, "replace", nil); err != nil {
		return fmt.Errorf("Failed starting service: %s, due to: %s\n", name, err)
	}

	return nil
}
