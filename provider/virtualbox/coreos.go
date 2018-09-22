package virtualbox

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	downloadURLFormat = "https://%s.release.core-os.net/amd64-usr/current/coreos_production_virtualbox_image.vmdk.bz2"
	versionURLFormat  = "https://%s.release.core-os.net/amd64-usr/current/version.txt"

	prefix           = "vmCoreOSImages"
	compressedName   = "coreos_production_virtualbox_image.vmdk.bz2"
	uncompressedName = "coreos_production_virtualbox_image.vmdk"
)

func BasePath(channel string, version string) string {
	return path.Join(os.TempDir(), prefix, channel, version)
}

// CompressedPath takes a channel and a version,
// and returns a path of where the compressed image should be stored.
func CompressedPath(channel string, version string) string {
	return path.Join(BasePath(channel, version), compressedName)
}

// CompressedPath takes a channel and a version,
// and returns a path of where the uncompressed image should be stored.
func UncompressedPath(channel string, version string) string {
	return path.Join(BasePath(channel, version), uncompressedName)
}

// LatestVersion takes a channel, and returns the latest version.
func LatestVersion(channel string) (string, error) {
	// TODO: Validate channel.

	versionURL := fmt.Sprintf(versionURLFormat, channel)
	resp, err := http.Get(versionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	versionPart := strings.Split(string(body), "\n")[3]
	// TODO: Validate versionPart.
	version := strings.Split(versionPart, "=")[1]
	// TODO: Validate version.

	return version, nil
}

// DownloadImage takes a channel and a path,
// and downloads the latest version of that channel to that path.
func DownloadImage(channel string, imagePath string) error {
	// TODO: Validate channel.

	dir := path.Dir(imagePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	downloadURL := fmt.Sprintf(downloadURLFormat, channel)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

// DecompressImage takes a path to an image, and decompresses it.
func DecompressImage(compressedImagePath string) error {
	// TODO: Can we do this without shelling out?
	return exec.Command("bunzip2", "--keep", compressedImagePath).Run()
}
