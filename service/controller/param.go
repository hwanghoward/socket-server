package controller

import (
	. "socket-server"
	"socket-server/service/models"
)

/*****************************Request param*************************************/

type CreateRaidParam struct {
	Name         []byte /////RAID名称
	Level        uint32 /////RAID级别
	Chunk        uint32 /////RAID条带大小    4、8、16、32、64、128、256、512、1024  单位为KB
	DataDisks    uint32
	SpareDisks   uint32
}

type HotSpareParam struct {
	RaidName     []byte ////Raid 唯一标识
	Location     uint32   ////槽位索引
}

type RaidIdentity struct {
	Name         []byte
}

type DiskIdentity struct {
	Location     uint32 ////槽位索引
}

type NoParam struct {

}

/*****************************New request param**********************************/

func NewCreateRaidParam() interface{} {
	return &CreateRaidParam{Name: make([]byte, 16)}
}

func NewHotSpareParam() interface{} {
	return &HotSpareParam{RaidName: make([]byte, 16)}
}

func NewRaidIdentity() interface{} {
	return &RaidIdentity{Name: make([]byte, 16)}
}

func NewDiskIdentity() interface{} {
	return &DiskIdentity{}
}

func NewNoParam() interface{} {
	return &NoParam{}
}

func NewParams(f func()interface{}, paramLens ...int) func(dataLen int) interface{} {
	return func(dataLen int) interface{} {
		found := false
		for _, paramLen := range paramLens {
			if paramLen == dataLen {
				found = true
				break
			}
		}
		if !found {
			return NewReplyFailResult(1, "Param invalid!")
		}

		if dataLen == 0 {
			return NewNoParam()
		}
		return f()
	}
}


/*****************************Result*********************************/

type RaidInfo struct {
	Name           []byte //16
	Level          uint32
	Chunk          uint32
	UUID           []byte//64
	TotalCap       uint32
	RemainCap      uint32
	State          uint32
}

type RaidRebuildInfo struct {
	Name           []byte //16
	UUID           []byte //64
	SyncPercent    uint32
}

type DiskInfo struct {
	SlotIndex      uint32
	Model          []byte //32
	Serial         []byte //64
	Fireware       []byte //32
	Capacity       uint64
	State          []byte //8
	Role           []byte //16
	RaidUUID       []byte //64
	DiskDev        []byte //16
}

type EnclosureInfo struct {
	ValidHDCount   uint32
	DiskBitarray   uint32
	RaidCount      uint32
}

type RaidInfoBytes struct {
	Info []byte    //108*n
}

type DiskInfoBytes struct {
	Info []byte    //240*n
}

func NewEnclosureInfo(hdCount, raidCount int, diskBitMap uint32) *EnclosureInfo {
	return &EnclosureInfo{ValidHDCount: uint32(hdCount), DiskBitarray: diskBitMap, RaidCount: uint32(raidCount)}
}

func NewRaidRebuildInfo(rd models.Raid) *RaidRebuildInfo {
	return &RaidRebuildInfo{
		Name:        string2Bytes(rd.Name, 16),
		UUID:        string2Bytes(rd.Uuid, 64),
		SyncPercent: uint32(rd.RbProgress),
	}
}