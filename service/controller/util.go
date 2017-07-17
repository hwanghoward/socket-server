package controller

import (
	"encoding/json"
	"fmt"
	"regexp"

	"socket-server"
	"socket-server/service/models"
	"strconv"
)

type ResErr struct {
	Status      string `json:"status"`
	ErrCode     int64  `json:"errcode"`
	Description string `json:"description"`
}

type ResRaid struct {
	Status string        `json:"status"`
	Detail []models.Raid `json:"detail"`
}

type ResDisk struct {
	Status string        `json:"status"`
	Detail []models.Disk `json:"detail"`
}

func bytes2String(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func bitmap2Slots(bitmap uint32) (slots []string) {
	for i := 0; i < 32; i++ {
		if bitmap&0x80000000 == uint32(0x80000000) {
			slots = append(slots, fmt.Sprintf("1.1.%d", i+1))
		}

		bitmap = bitmap << 1
	}

	return
}

func slots2Bitmap(slots []string) (bitmap uint32) {
	for _, slot := range slots {
		re := regexp.MustCompile(`\d\.\d\.(\d+)`)
		res := re.FindStringSubmatch(slot)
		if len(res) == 2 {
			i, _ := strconv.Atoi(res[1])
			i = i - 1
			flag := uint32(0x80000000)
			bitmap |= (flag >> uint(i))
		}
	}

	return
}

func slot2Index(slot string) (index uint32) {
	re := regexp.MustCompile(`\d\.\d\.(\d+)`)
	res := re.FindStringSubmatch(slot)
	if len(res) == 2 {
		i, _ := strconv.Atoi(res[1])
		i = i - 1
		flag := uint32(0x80000000)
		index |= (flag >> uint(i))
	}

	return
}

func string2Bytes(str string, length int) []byte {
	bs := []byte(str)
	if len(bs) < length {
		add := make([]byte, length-len(bs))
		bs = append(bs, add...)
		return bs
	}
	return bs[0 : length]
}

func state2Uint32(stat string, rebuilding bool) uint32 {
	if rebuilding {
		return uint32(RAID_STAT["rebuilding"])
	}
	return uint32(RAID_STAT[stat])
}

func parseError(res []byte) (result ResErr) {
	if err := json.Unmarshal(res, &result); err != nil {
	}
	return
}

func raidInfo2Bytes(rdGoInfos []models.Raid) *RaidInfoBytes {
	rdCInfos := make([]*RaidInfo, 0)
	for _, rdGoInfo := range rdGoInfos {
		rdCInfo := &RaidInfo{
			Name:      string2Bytes(rdGoInfo.Name, 16),
			Level:     uint32(RAID_LEVEL2[int(rdGoInfo.Level)]),
			Chunk:     uint32(rdGoInfo.ChunkKb),
			UUID:      string2Bytes(rdGoInfo.Uuid, 64),
			TotalCap:  uint32(rdGoInfo.CapMb),
			RemainCap: uint32(rdGoInfo.CapMb - rdGoInfo.UsedMb),
			State:     state2Uint32(rdGoInfo.Health, rdGoInfo.Rb),
		}
		rdCInfos = append(rdCInfos, rdCInfo)
	}

	infoBytes := &RaidInfoBytes{
		Info : make([]byte, 0),
	}
	for _, rdCInfo := range rdCInfos {
		bs := socket.MarshalParam(rdCInfo)
		infoBytes.Info = append(infoBytes.Info, bs...)
	}

	return infoBytes
}

func diskInfo2Bytes(hdGoInfos []models.Disk) *DiskInfoBytes {
	hdCInfos := make([]*DiskInfo, 0)

	for _, hdGoInfo := range hdGoInfos {
		hdCInfo := &DiskInfo{
			SlotIndex: slot2Index(hdGoInfo.Location),
			Model:     string2Bytes(hdGoInfo.Uuid, 64),
			Serial:    string2Bytes(hdGoInfo.Sn, 64),
			Fireware:  string2Bytes(hdGoInfo.Vendor, 32),
			Capacity:  uint64(hdGoInfo.CapSector * 512),
			State:     string2Bytes(hdGoInfo.Host, 8),
			Role:      string2Bytes(hdGoInfo.Role, 16),
			RaidUUID:  string2Bytes(hdGoInfo.Raid, 64),
			DiskDev:   string2Bytes(hdGoInfo.DevName, 16),
		}
		hdCInfos = append(hdCInfos, hdCInfo)
	}

	infoBytes := &DiskInfoBytes{
		Info: make([]byte, 0),
	}
	for _, hdCInfo := range hdCInfos {
		bs := socket.MarshalParam(hdCInfo)
		infoBytes.Info = append(infoBytes.Info, bs...)
	}

	return infoBytes
}

