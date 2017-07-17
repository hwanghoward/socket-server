package models

import (
	"fmt"
	"regexp"
	"strings"

	"socket-server/service/util"
)

type Raid struct {
	Rb         bool    `json:"rebuilding"`
	Uuid       string  `json:"id"`
	Health     string  `json:"health"`
	Level      int64   `json:"level"`
	Name       string  `json:"name"`
	Cap        int64   `json:"cap_sector"`
	Used       int64   `json:"used_cap_sector"`
	CapMb      float64 `json:"cap_mb"`
	UsedMb     float64 `json:"used_cap_mb"`
	ChunkKb    int64   `json:"chunk_kb"`
	Blkdev     string  `json:"blkdev"`
	RbProgress float64 `json:"rebuild_progress"`
}

func GetAllRaids(uuid string) (r string, err error) {
	cmd := fmt.Sprintf("python speedcmd.py --method=list_raids --argument='uuid=%s'", uuid)
	raid, err := util.ExecuteByStr(cmd, true)
	if err != nil {

	}

	m := regexp.MustCompile(`\{.*\}`)
	raids := m.FindStringSubmatch(raid)
	if len(raids) > 0 {
		r = raids[0]
		return
	}
	return
}

func AddRaids(name, level, chunk, raid_disks, spare_disks, rebuild_priority, sync string) (result string, err error) {
	cmd := fmt.Sprintf("python speedcmd.py --method=create_raid --argument='name=%s&level=%s&chunk=%s&raid_disks=%s&spare_disks=%s&rebuild_priority=%s&sync=%s'",
		name, level, chunk, raid_disks, spare_disks, rebuild_priority, sync)

//	cmd := fmt.Sprintf("python speedcmd.py --method=create_raid --argument='name=rd&level=%s&chunk=%s&raid_disks=%s&spare_disks=%s&rebuild_priority=%s&sync=%s'",
//		level, chunk, raid_disks, spare_disks, rebuild_priority, sync)

	re, err := util.ExecuteByStr(cmd, true)
	if err != nil {
	}

	m := regexp.MustCompile(`\{.*\}`)
	res := m.FindStringSubmatch(re)
	if len(res) > 0 {
		result = res[0]
		if strings.Contains(result, "error") {
			err = fmt.Errorf("error")
		}
		return
	}
	return
}

func DelRaids(name string) (result string, err error) {
	cmd := fmt.Sprintf("python speedcmd.py --method=delete_raid --argument='name=%s'", name)
	re, err := util.ExecuteByStr(cmd, true)
	if err != nil {
	}

	m := regexp.MustCompile(`\{.*\}`)
	res := m.FindStringSubmatch(re)
	if len(res) > 0 {
		result = res[0]
		if strings.Contains(result, "error") {
			err = fmt.Errorf("error")
		}
		return
	}

	return
}
