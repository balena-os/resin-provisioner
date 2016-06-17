package provisioner

import (
	"fmt"
	"log"
)

func getOsReleaseSlug() (string, error) {
	if fields, err := getEnvFileFields(OSRELEASE_PATH); err != nil {
		return "", err
	} else {
		slug := fields["SLUG"]
		if slug == "" {
			return "", fmt.Errorf("Could not find 'SLUG' field in %s",
				OSRELEASE_PATH)
		}

		return slug, nil
	}
}

func ScanDeviceTypeSlug() (string, error) {
	// We first attempt to retrieve the slug from /etc/os-release.
	if slug, err := getOsReleaseSlug(); err == nil {
		return slug, nil
	} else {
		log.Printf("WARNING: Could not retrieve os-release slug, "+
			"reverting to heuristics: %s\n", err)

		deviceType, err := identify()
		return deviceType.String(), err
	}
}
