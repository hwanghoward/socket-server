package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	. "socket-server"
	"socket-server/service/models"
)

type Operator struct {
}

func (this *Operator) Excute(message *Message) {
	switch message.Cmd {
	case CMD_REQUEST_CREATE_RAID:
		this.createRaid(message)
	case CMD_REQUEST_DETORY_RAID:
		this.deleteRaid(message)
	case CMD_REQUEST_CREATE_HOTSPARE:
		this.createHotSpare(message)
	case CMD_REQUEST_DETORY_HOTSPARE:
		this.destroyHotSpare(message)
	case CMD_REQUEST_FORMAT_DISK:
		this.formatDisk(message)
	default:
	}

	return
}

func (this *Operator) createRaid(message *Message) {
	param, ok := message.Param.(*CreateRaidParam)
	if !ok {
		return
	}

	name := bytes2String(param.Name)
	level := RAID_LEVEL[int(param.Level)]
	chunk := fmt.Sprintf("%dKB", int(param.Chunk))
	raidDisks := strings.Join(bitmap2Slots(param.DataDisks), ",")
	spareDisks := strings.Join(bitmap2Slots(param.SpareDisks), ",")
	sync := "no"
	rebuildPriority := "low"

	raidRe, err := models.AddRaids(name, level, chunk, raidDisks, spareDisks, rebuildPriority, sync)
	if err != nil {
		res := parseError([]byte(raidRe))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_CREATE_RAID_FAIL)
		message.Param = param
		return
	}

	var resRaid ResRaid
	resRaid.Detail = make([]models.Raid, 0)
	if err = json.Unmarshal([]byte(raidRe), &resRaid); err != nil {
	}

	reParam := NewNoParam()
	message.Cmd = uint32(CMD_REPLY_CREATE_RAID_SUCCESS)
	message.Param = reParam
	return
}

// delete raid with uuid
func (this *Operator) deleteRaid(message *Message) {
	param, ok := message.Param.(*RaidIdentity)
	if !ok {
		return
	}

	name := bytes2String(param.Name)
	raidRe, err := models.DelRaids(name)
	if err != nil {
		res := parseError([]byte(raidRe))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_DETORY_RAID_FAIL)
		message.Param = param
		return
	}

	reParam := NewNoParam()
	message.Cmd = uint32(CMD_REPLY_DETORY_RAID_SUCCESS)
	message.Param = reParam
	return
}

// create hot spare
func (this *Operator) createHotSpare(message *Message) {
	param, ok := message.Param.(*HotSpareParam)
	if !ok {
		return
	}

	name := bytes2String(param.RaidName)
	slots := bitmap2Slots(param.Location)
	if len(slots) == 0 {
		return
	}
	location := slots[0]

	raidRe, err := models.UpdateDisks(location, "", "spare", name)
	if err != nil {
		res := parseError([]byte(raidRe))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_CREATE_HOTSPARE_FAIL)
		message.Param = param
		return
	}

	reParam := NewNoParam()
	message.Cmd = uint32(CMD_REPLY_CREATE_HOTSPARE_SUCCESS)
	message.Param = reParam
	return
}

func (this *Operator) destroyHotSpare(message *Message) {
	param, ok := message.Param.(*HotSpareParam)
	if !ok {
		return
	}

	name := bytes2String(param.RaidName)
	slots := bitmap2Slots(param.Location)
	if len(slots) == 0 {
		return
	}

	location := slots[0]

	raidRe, err := models.UpdateDisks(location, "", "unused", name)
	if err != nil {
		res := parseError([]byte(raidRe))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_DETORY_HOTSPARE_FAIL)
		message.Param = param
		return
	}

	reParam := NewNoParam()
	message.Cmd = uint32(CMD_REPLY_DETORY_HOTSPARE_SUCCESS)
	message.Param = reParam
	return
}

// format disk
func (this *Operator) formatDisk(message *Message) {
	param, ok := message.Param.(*DiskIdentity)
	if !ok {
		return
	}

	slots := bitmap2Slots(param.Location)
	if len(slots) == 0 {
		return
	}

	location := slots[0]
	raidRe, err := models.UpdateDisks(location, "native", "", "")
	if err != nil {
		res := parseError([]byte(raidRe))
		param := NewReplyFailResult(uint32(res.ErrCode), res.Description)
		message.Cmd = uint32(CMD_REPLY_FORMAT_DISK_FAIL)
		message.Param = param
		return
	}

	reParam := NewNoParam()
	message.Cmd = uint32(CMD_REPLY_FORMAT_DISK_SUCCESS)
	message.Param = reParam
	return
}
