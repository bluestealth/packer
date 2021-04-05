package firmware

import (
	"strings"
)

type (
	FirmwareOSInterface string
	FirmwareDevice      string
)

const (
	BIOS         FirmwareOSInterface = "bios"
	OPENFIRMWARE                     = "openfirmware"
	UBOOT                            = "uboot"
	UEFI                             = "uefi"
)

const (
	FLASH  FirmwareDevice = "flash"
	KERNEL                = "kernel"
	MEMORY                = "memory"
)

type FirmwareTarget struct {
	Architecture string   `json:"architecture"`
	Machines     []string `json:"machines"`
}

type FirmwareFlashFile struct {
	Filename string `json:"filename"`
	Format   string `json:"format"`
}

type FirmwareMapping struct {
	Device           FirmwareDevice     `json:"device"`
	NonFlashFilename string             `json:"filename,omitempty"`
	Executable       *FirmwareFlashFile `json:"executable,omitempty"`
	NVRAMTemplate    *FirmwareFlashFile `json:"nvram-template,omitempty"`
}

type Firmware struct {
	Priority       int                   `json:"-"`
	Filename       string                `json:"-"`
	Description    string                `json:"description"`
	InterfaceTypes []FirmwareOSInterface `json:"interface-types"`
	Mapping        FirmwareMapping       `json:"mapping"`
	Targets        []FirmwareTarget      `json:"targets"`
	Features       []string              `json:"features"`
	Tags           []string              `json:"tags"`
}

func (f Firmware) TargetMatches(architecture, machine string) bool {
	switch machine {
	case "pc":
		machine = "pc-i440fx"
	case "q35":
		machine = "pc-q35"
	}

	for _, target := range f.Targets {
		if target.Architecture == architecture {
			for _, targetMachine := range target.Machines {
				if strings.HasPrefix(machine, strings.TrimSuffix(targetMachine, "-*")) {
					return true
				}
			}
		}
	}
	return false
}

func (f Firmware) HasFeature(target string) bool {
	for _, feature := range f.Features {
		if feature == target {
			return true
		}
	}
	return false
}

func (f Firmware) HasInterface(target FirmwareOSInterface) bool {
	for _, interfaceType := range f.InterfaceTypes {
		if interfaceType == target {
			return true
		}
	}
	return false
}
