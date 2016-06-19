package provisioner

import "log"

func getOsReleaseSlug() (string, error) {
	if fields, err := getEnvFileFields(OSRELEASE_PATH); err != nil {
		return "", err
	} else {
		slug := fields["SLUG"]
		if slug == "" {
			log.Printf("Could not find 'SLUG' field in %s, defaulting to raspberry-pi\n",
				OSRELEASE_PATH)
			slug = "raspberry-pi"
		}

		return slug, nil
	}
}

func ScanDeviceTypeSlug() (string, error) {
	return getOsReleaseSlug()
}
