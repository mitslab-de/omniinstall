package platform

import (
	"io"
	"os"
)

// openOsRelease opens /etc/os-release for reading.
func openOsRelease() (io.ReadCloser, error) {
	return os.Open("/etc/os-release")
}
