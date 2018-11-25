package protocol

import (
	"bytes"
	"encoding/binary"
)

//convert.go文件定义了字节流与对应结构体的相互转换方法
//虽然有一个unsafe包好像可以直接转，不过这个包名字起的我就不敢用..以后可以研究一下怎么用

func BytesToInt(b []byte, i interface{}) interface{} {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.LittleEndian, i)
	return i
}

func IntToBytes(n interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

func convertBytesToHeader(headBuffer []byte) *Head {
	head 		:= new(Head)
	headSize 	:= HeadSize
	if uint64(len(headBuffer)) < headSize {
		return nil
	}

	head.Mark = string(headBuffer[0:4])
	BytesToInt(headBuffer[4:12], &(head.BodyLength))
	copy(head.BodyMd5[0:16], headBuffer[12:28])
	copy(head.HeadMd5[0:16], headBuffer[28:44])
	return head
}

