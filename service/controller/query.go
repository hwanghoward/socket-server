package controller

import (
//	"fmt"
	"encoding/json"

	. "socket-server"
	"socket-server/service/models"
)

type Query struct {
}

func (this *Query) Excute(message *Message) {
	switch message.Cmd {
	case CMD_REQUEST_QUERY_ENCLORSURE_INFO:
		this.queryEnclosure(message)
	case CMD_REQUEST_QUERY_RAID_INFO:
		this.queryRaid(message)
	case CMD_REQUEST_QUERY_DISK_INFO:
		this.queryDisk(message)
	case CMD_REQUEST_QUERY_RAID_REBUILD_INFO:
		this.queryRaidRebuild(message)
	default:
	}

	return
}

// get enclosure info
func (this *Query) queryEnclosure(message *Message) {
	_, ok := message.Param.(*NoParam)
	if !ok {
		return
	}

	//get raids
	raids, err := models.GetAllRaids("")
	if err != nil {
		res := parseError([]byte(raids))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_QUERY_ENCLORSURE_INFO_FAIL)
		message.Param = param
		return
	}
	var resRaid ResRaid
	resRaid.Detail = make([]models.Raid, 0)
	if err = json.Unmarshal([]byte(raids), &resRaid); err != nil {
	}

	//get disks
	disks, err := models.GetRestDisks()
	if err != nil {
		res := parseError([]byte(raids))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_QUERY_ENCLORSURE_INFO_FAIL)
		message.Param = param
		return
	}
	var resDisk ResDisk
	resDisk.Detail = make([]models.Disk, 0)
	if err = json.Unmarshal([]byte(disks), &resDisk); err != nil {
	}

	//get bitmap
	diskLocations := make([]string, 0)
	for _, hd := range resDisk.Detail {
		diskLocations = append(diskLocations, hd.Location)
	}
	bitmap := slots2Bitmap(diskLocations)

	reParam := NewEnclosureInfo(len(resDisk.Detail), len(resRaid.Detail), bitmap)
	message.Cmd = uint32(CMD_REPLY_QUERY_ENCLORSURE_INFO_SUCCESS)
	message.Param = reParam
	return
}

func (this *Query) queryRaid(message *Message) {
	name := ""
	if param, ok := message.Param.(*RaidIdentity); ok {
		name = bytes2String(param.Name)
	}

	//get raids
	raids, err := models.GetAllRaids(name)
	if err != nil {
		res := parseError([]byte(raids))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_QUERY_RAID_INFO_FAIL)
		message.Param = param
		return
	}
	var resRaid ResRaid
	resRaid.Detail = make([]models.Raid, 0)
	if err = json.Unmarshal([]byte(raids), &resRaid); err != nil {
	}

	reParam := raidInfo2Bytes(resRaid.Detail)
	message.Cmd = uint32(CMD_REPLY_QUERY_RAID_INFO_SUCCESS)
	message.Param = reParam
	return
}

func (this *Query) queryDisk(message *Message) {
	_, ok := message.Param.(*NoParam)
	if !ok {
		return
	}

	//get disks
	disks, err := models.GetRestDisks()
	if err != nil {
		res := parseError([]byte(disks))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_QUERY_DISK_INFO_FAIL)
		message.Param = param
		return
	}
	var resDisk ResDisk
	//resDisk.Detail = make([]models.Disk, 0)
	if err = json.Unmarshal([]byte(disks), &resDisk); err != nil {
	}
	//fmt.Println(resDisk)

	reParam := diskInfo2Bytes(resDisk.Detail)
	message.Cmd = uint32(CMD_REPLY_QUERY_DISK_INFO_SUCCESS)
	message.Param = reParam
	return
}

func (this *Query) queryRaidRebuild(message *Message) {
	param, ok := message.Param.(*RaidIdentity)
	if !ok {
		return
	}

	name := bytes2String(param.Name)
	//get raids
	raids, err := models.GetAllRaids(name)
	if err != nil {
		res := parseError([]byte(raids))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_QUERY_RAID_REBUILD_INFO_FAIL)
		message.Param = param
		return
	}
	var resRaid ResRaid
	resRaid.Detail = make([]models.Raid, 0)
	if err = json.Unmarshal([]byte(raids), &resRaid); err != nil {
	}

	if len(resRaid.Detail) == 0 {
		param := NewReplyFailResult(uint32(2), "raid no found!")
		message.Cmd = uint32(CMD_REPLY_QUERY_RAID_REBUILD_INFO_FAIL)
		message.Param = param
		return
	}

	reParam := NewRaidRebuildInfo(resRaid.Detail[0])
	message.Cmd = uint32(CMD_REPLY_QUERY_RAID_REBUILD_INFO_SUCCESS)
	message.Param = reParam
	return
}
