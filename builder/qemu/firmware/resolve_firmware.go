package firmware

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func getLocation() string {
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "solaris":
		return "/usr/share/qemu/firmware"
	case "darwin":
		return "/usr/local/share/qemu/firmware"
	default:
		return ""
	}
}

func GetFirmwares() (results []*Firmware, err error) {
	firmwareDir := getLocation()
	if firmwareDir == "" {
		return
	}

	if _, err = os.Stat(firmwareDir); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}

	err = filepath.Walk(firmwareDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() {
			basename := filepath.Base(path)
			if ok, _ := filepath.Match("*.json", basename); ok {
				priority := strings.SplitN(basename, "-", 2)[0]
				priorityVal, err := strconv.Atoi(priority)
				if err != nil {
					return err
				}
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				var firmwareInfo *Firmware
				if err := json.Unmarshal(data, &firmwareInfo); err != nil {
					return err
				}
				firmwareInfo.Filename = path
				firmwareInfo.Priority = priorityVal
				results = append(results, firmwareInfo)
			}
		}
		return nil
	})

	return
}
