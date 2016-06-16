package provisioner

// A super heuristic set of identification criteria. Forgive me, Linus.

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	CPUINFO_PATH      = "/proc/cpuinfo"
	PRODUCT_NAME_PATH = "/sys/devices/virtual/dmi/id/product_name"
	BOARD_NAME_PATH   = "/sys/devices/virtual/dmi/id/board_name"
)

type DeviceType int

const (
	UnknownDevice DeviceType = iota
	RaspberryPi0
	RaspberryPi1
	RaspberryPi2
	RaspberryPi3
	Edison
	Nuc
	BeagleboneBlack
)

func (dt DeviceType) String() string {
	switch dt {
	case RaspberryPi0:
		fallthrough
	case RaspberryPi1:
		return "raspberry-pi"
	case RaspberryPi2:
		return "raspberry-pi2"
	case RaspberryPi3:
		return "raspberrypi3"
	case Edison:
		return "intel-edison"
	case Nuc:
		return "intel-nuc"
	case BeagleboneBlack:
		return "beaglebone-black"
	}

	return "???"
}

func readFile(path string) (string, error) {
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return "", fmt.Errorf("Can't open %s: %s.", PRODUCT_NAME_PATH, err)
	} else {
		return strings.TrimSpace(string(bytes)), nil
	}
}

func parseProcCpu(procCpu string) (hardware, revision string) {
	getVal := func(line string) string {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			return ""
		}

		return strings.TrimSpace(fields[1])
	}

	for _, line := range strings.Split(procCpu, "\n") {
		if strings.HasPrefix(line, "Hardware") {
			hardware = getVal(line)
		} else if strings.HasPrefix(line, "Revision") {
			revision = getVal(line)
		}
	}

	return
}

func identifyBeaglebone(hardware string) (DeviceType, error) {
	if strings.HasSuffix(hardware, "AM33XX") {
		return BeagleboneBlack, nil
	}

	return UnknownDevice, fmt.Errorf("Unrecognised hardware %s.", hardware)
}

func identifyRpi(revision string) (ret DeviceType, err error) {
	if revision == "" {
		return UnknownDevice, fmt.Errorf("Couldn't find revision field.")
	}

	// Revision is a hex number.
	// Cribbed revisions from http://elinux.org/RPi_HardwareHistory
	if n, err := strconv.ParseInt(revision, 16, 32); err != nil {
		return UnknownDevice, fmt.Errorf("Couldn't parse rpi revision '%s'.",
			revision)
	} else if n <= 0x15 {
		return RaspberryPi1, nil
	} else if n == 0x900092 || n == 0x900093 {
		return RaspberryPi0, nil
	} else if n == 0xa01041 || n == 0xa21041 {
		return RaspberryPi2, nil
	} else if n == 0xa02082 || n == 0xa22082 {
		return RaspberryPi3, nil
	}

	return UnknownDevice, nil
}

func identifyArm() (DeviceType, error) {
	var (
		procCpu string
		err     error
	)

	if procCpu, err = readFile(CPUINFO_PATH); err != nil {
		return UnknownDevice, err
	}

	hardware, revision := parseProcCpu(procCpu)

	if hardware == "" {
		return UnknownDevice, fmt.Errorf("Couldn't find hardware field in %s.",
			CPUINFO_PATH)
	}

	if strings.HasPrefix(hardware, "BCM") {
		return identifyRpi(revision)
	}

	return identifyBeaglebone(hardware)
}

func identifyIntel() (DeviceType, error) {
	var (
		productName, boardName string
		err                    error
	)

	if productName, err = readFile(PRODUCT_NAME_PATH); err != nil {
		return UnknownDevice, err
	}

	if strings.HasPrefix(productName, "NU") {
		return Nuc, nil
	}

	if productName == "Merrifield" {
		return Edison, nil
	}

	// Next try board name.
	if boardName, err = readFile(BOARD_NAME_PATH); err != nil {
		return UnknownDevice, err
	}

	if strings.HasPrefix(boardName, "NU") {
		return Nuc, nil
	}
	if strings.HasPrefix(boardName, "D") {
		return Nuc, nil
	}

	// Nuc is a pretty good bet.
	fmt.Fprintf(os.Stderr,
		"WARNING: Can't identify x86 device via product/board name %s/%s, "+
			"assuming NUC.\n",
		productName, boardName)
	return Nuc, nil
}

func ScanDeviceType() (DeviceType, error) {
	if runtime.GOARCH == "arm" {
		return identifyArm()
	}

	return identifyIntel()
}
