package models

import (
	_ "encoding/json"
	"fmt"
	"regexp"
	"strings"

	"socket-server/service/util"
)

type Disk struct {
	Uuid      string  `json:"id"`
	Host      string  `json:"host"`
	Health    string  `json:"health"`
	Role      string  `json:"role"`
	Location  string  `json:"location"`
	Raid      string  `json:"raid"`
	CapSector int64   `json:"cap_sector"`
	CapMb     float64 `json:"cap_mb"`
	Vendor    string  `json:"vendor"`
	Model     string  `json:"model"`
	Sn        string  `json:"sn"`
	RqrCount  int64   `json:"rqr_count"`
	Rpm       int64   `json:"rpm"`
	DevName   string  `json:"dev_name"`
}

func GetRestDisks() (d string, err error) {
	cmd := "python speedcmd.py"
	disk, err := util.ExecuteByStr(cmd, true)
	if err != nil {
	}

	m := regexp.MustCompile(`\{.*\}`)
	disks := m.FindStringSubmatch(disk)
	if len(disks) > 0 {
		d = disks[0]
		return
	}
	return
}

func UpdateDisks(location, host, role, raid string) (result string, err error) {
	var o string
	if host == "native" {
		cmd := fmt.Sprintf("python speedcmd.py --method=update_disk --argument='location=%s&host=native'", location)
		//	cmd = fmt.Sprintf("python speedcmd.py --method=update_disk --argument=location=%s&host=%s&role=%s&raid=%s", d.Location, d.Host, d.Role, d.Raid)
		o, err = util.ExecuteByStr(cmd, true)
		if err != nil {
		}
	} else {
		cmd := fmt.Sprintf("python speedcmd.py --method=update_disk --argument='location=%s&role=%s&raid=%s'", location, role, raid)
		o, err = util.ExecuteByStr(cmd, true)
		if err != nil {
		}

	}
	m := regexp.MustCompile(`\{.*\}`)
	disks := m.FindStringSubmatch(o)
	if len(disks) > 0 {
		result = disks[0]
		if strings.Contains(result, "error") {
			err = fmt.Errorf("error")
		}
	}

	return
}
