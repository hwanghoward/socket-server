package controller

const (
	CMD_BASENUM_RAID_OPERATOR       = 0x200
	CMD_REQUEST_CREATE_RAID         = iota + CMD_BASENUM_RAID_OPERATOR
	CMD_REQUEST_DETORY_RAID
	CMD_REQUEST_CREATE_HOTSPARE
	CMD_REQUEST_DETORY_HOTSPARE
	CMD_REQUEST_FORMAT_DISK
)

const (
	CMD_BASENUM_REPLY_RAID_OPERATOR = 0x1000200
	CMD_REPLY_CREATE_RAID_SUCCESS   = iota + CMD_BASENUM_REPLY_RAID_OPERATOR
	CMD_REPLY_CREATE_RAID_FAIL
	CMD_REPLY_DETORY_RAID_SUCCESS
	CMD_REPLY_DETORY_RAID_FAIL
	CMD_REPLY_CREATE_HOTSPARE_SUCCESS
	CMD_REPLY_CREATE_HOTSPARE_FAIL
	CMD_REPLY_DETORY_HOTSPARE_SUCCESS
	CMD_REPLY_DETORY_HOTSPARE_FAIL
	CMD_REPLY_FORMAT_DISK_SUCCESS
	CMD_REPLY_FORMAT_DISK_FAIL
)

const (
	CMD_BASENUM_QUERY                    = 0x300
	CMD_REQUEST_QUERY_ENCLORSURE_INFO    = iota + CMD_BASENUM_QUERY
	CMD_REQUEST_QUERY_RAID_INFO
	CMD_REQUEST_QUERY_DISK_INFO
	CMD_REQUEST_QUERY_RAID_REBUILD_INFO
)

const (
	CMD_BASENUM_REPLY_QUERY                  = 0x1000300
	CMD_REPLY_QUERY_ENCLORSURE_INFO_SUCCESS  = iota + CMD_BASENUM_REPLY_QUERY
	CMD_REPLY_QUERY_ENCLORSURE_INFO_FAIL
	CMD_REPLY_QUERY_RAID_INFO_SUCCESS
	CMD_REPLY_QUERY_RAID_INFO_FAIL
	CMD_REPLY_QUERY_DISK_INFO_SUCCESS
	CMD_REPLY_QUERY_DISK_INFO_FAIL
	CMD_REPLY_QUERY_RAID_REBUILD_INFO_SUCCESS
	CMD_REPLY_QUERY_RAID_REBUILD_INFO_FAIL
)

const (
	RAID_0 = iota + 1
	RAID_1
	RAID_5
	RAID_6
	RAID_10
)

var RAID_LEVEL = map[int]string {
	RAID_0:  "0",
	RAID_1:  "1",
	RAID_5:  "5",
	RAID_6:  "6",
	RAID_10: "10",
}

var RAID_LEVEL2 = map[int]int {
	0: RAID_0,
	1: RAID_1,
	5: RAID_5,
	6: RAID_6,
	10: RAID_10,
}

const (
    NORMAL = iota + 1////正常
    DEGRADE////降级
    FAULT////故障
    REBUILD////重建
    INITIAL////初始化中
)

var RAID_STAT = map[string]int {
	"normal":    NORMAL,
	"degraded":  DEGRADE,
	"failed":    FAULT,
	"rebuilding":REBUILD,
	"initial":   INITIAL,
}
