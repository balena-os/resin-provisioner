package provisioner

import "fmt"

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
	return getOsReleaseSlug()
}
