//通讯协议处理，主要处理封包和解包的过程
package socket

import (
	"fmt"
	"reflect"
	"bytes"
	"encoding/binary"
)

const (
	constHeader       = "zyrd"

	constHeaderLength = 4
	versionLength     = 4
	idLength          = 4
	userPragramLength = 4
	dataLenLength     = 4
)

//Header
type PackCmdHeader struct {
	CmdHeader   []byte  //zyrd
	CmdVersion  uint32  //协议版本
	CmdID       uint32  //命令ID
	UserParam   uint32  //答复时携带，内容相同
	DataLen     uint32  //除本结构体外数据长度
}

//Error type
type ReplyFailResult struct {
	ErrorNo        uint32
	DescriptionLen uint32
	Description    []byte //126
}


//封包
func packet(msg *Message) []byte {
	headerObj := PackCmdHeader{
		CmdHeader:  []byte(constHeader),
		CmdVersion: msg.Version,
		CmdID:      msg.Cmd,
		UserParam:  msg.UserProgram,
	}
	headerBytes := MarshalParam(&headerObj)

	data := MarshalParam(msg.Param)
	return append(headerBytes, data...)
}

//解包
func unpack(buffer []byte, readerChannel chan *Message) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+constHeaderLength+versionLength+idLength+userPragramLength+dataLenLength {
			break
		}
		pos := i+constHeaderLength
		if string(buffer[i:pos]) == constHeader {
			version := BytesToUint32(buffer[pos : pos+versionLength])
			pos += versionLength
			id := BytesToUint32(buffer[pos : pos+idLength])
			pos += idLength
			userProgram := BytesToUint32(buffer[pos : pos+userPragramLength])
			pos += userPragramLength
			dataLen := BytesToUint32(buffer[pos : pos+dataLenLength])
			pos += dataLenLength

			if length < pos+int(dataLen) {
				break
			}
			data := buffer[pos : pos+int(dataLen)]

			msg := UnmarshalParam(id, version, userProgram, data)

			readerChannel <- msg
			i = pos+int(dataLen) - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func getParamStruct(cmd, dataLen int) interface{} {
	f, ok := Api[cmd]
	if !ok {
		return NewReplyFailResult(111, "Request no found!")
	}
	
	return f(dataLen)
}

func UnmarshalParam(cmd uint32, version uint32, userProgram uint32, data []byte) *Message {
	structT := getParamStruct(int(cmd), len(data))
	if _, ok := structT.(*ReplyFailResult); ok {
		return &Message{Cmd: cmd, Version: version, UserProgram: userProgram, Param: structT}
	}
	fmt.Printf("message: %v, %#v\n", structT, cmd)

	t := reflect.ValueOf(structT)
	if t.Kind() != reflect.Ptr {
		fmt.Println("Param structT should be a Ptr!")
		return nil
	}

	t = t.Elem()
	size := 0
	for i:=0; i < t.NumField(); i++ {
		field := t.Field(i)

		switch field.Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				leng := field.Type().Bits()/8
				val := BytesToUint64(data[size : size+leng])
				size += leng
				if field.CanSet() {
					field.Set(reflect.ValueOf(val).Convert(field.Type()))
				}
			case reflect.Slice:
				leng := field.Len()
				val := data[size : size+leng]
				size += leng
				if field.CanSet() {
					field.Set(reflect.ValueOf(val).Convert(field.Type()))
				}
		}
	}

	return &Message{Cmd: cmd, Version: version, UserProgram: userProgram, Param: structT}
}

func MarshalParam(structT interface{}) []byte {
	result := make([]byte, 0)

	t := reflect.ValueOf(structT)
	if t.Kind() != reflect.Ptr {
		fmt.Println("Param structT should be a Ptr!")
		return result
	}

	t = t.Elem()
	
	for i:=0; i < t.NumField(); i++ {
		field := t.Field(i)

		switch field.Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val := field.Uint()
				leng := field.Type().Bits()/8
				result = append(result, UintToBytes(val, leng)...)
			case reflect.Slice:
				result = append(result, field.Bytes()...)
		}
	}

	return result
}


/*整形转换成字节*/
func Uint8ToBytes(n uint8) (bs []byte) {
	return UintToBytes(uint64(n), 1)
}

func Uint16ToBytes(n uint16) (bs []byte) {
	return UintToBytes(uint64(n), 2)
}

func Uint32ToBytes(n uint32) (bs []byte) {
	return UintToBytes(uint64(n), 4)
}

func Uint64ToBytes(n uint64) (bs []byte) {
	return UintToBytes(n, 8)
}

func UintToBytes(n uint64, length int) (bs []byte) {
	if length > 8 || length <= 0 || length%2 != 0 {
		return 
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	bs = bytesBuffer.Bytes()
	return bs[8-length:]
}

/*字节转换成整形*/
func BytesToUint8(b []byte) uint8 {
	return uint8(BytesToUint64(b))
}

func BytesToUint16(b []byte) uint16 {
	return uint16(BytesToUint64(b))
}

func BytesToUint32(b []byte) uint32 {
	return uint32(BytesToUint64(b))
}

func BytesToUint64(b []byte) uint64 {
	var pre []byte
	if len(b) > 8 {
		pre = b[len(b)-8:]
	} else {
		pre = make([]byte, 8-len(b))
		pre = append(pre, b...)
	}

	bytesBuffer := bytes.NewBuffer(pre)

	var x uint64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

func NewReplyFailResult(errorNo uint32, des string) *ReplyFailResult {
	des2bytes := []byte(des)
	Len := len(des2bytes)
	if Len > 126 {
		return &ReplyFailResult{ErrorNo: errorNo, DescriptionLen: uint32(Len), Description: des2bytes[0:126]}
	}
	suffix := make([]byte, 126-Len)
	des2bytes = append(des2bytes, suffix...)
	return &ReplyFailResult{ErrorNo: errorNo, DescriptionLen: uint32(Len), Description: des2bytes}
}