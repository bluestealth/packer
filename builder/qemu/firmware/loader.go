package firmware

import (
	"sort"
)

func Loader(architecture, firmwareType, machineType string) (firmware *Firmware, err error) {
	var firmwares []*Firmware
	if firmwares, err = GetFirmwares(); err != nil {
		return
	}

	// Sort results by priority
	sort.SliceStable(firmwares, func(i, j int) bool {
		return firmwares[i].Priority < firmwares[j].Priority
	})

	switch firmwareType {
	case "bios":
		for _, fw := range firmwares {
			if fw.TargetMatches(architecture, machineType) &&
				fw.HasInterface(BIOS) {
				firmware = fw
				return
			}
		}
	case "uefi":
		for _, fw := range firmwares {
			if fw.TargetMatches(architecture, machineType) &&
				fw.HasInterface(UEFI) &&
				!fw.HasFeature("secure-boot") {
				firmware = fw
				return
			}
		}
	case "uefi-secure":
		for _, fw := range firmwares {
			if fw.TargetMatches(architecture, machineType) &&
				fw.HasInterface(UEFI) &&
				fw.HasFeature("secure-boot") {
				firmware = fw
				return
			}
		}
	}

	return
}
