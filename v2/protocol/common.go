package protocol

import (
	"errors"
	"net"
	"crypto/md5"
	"fmt"
	"reflect"
)

var HeadError = errors.New("获取头部信息出错")
var BodyError = errors.New("获取Body信息出错")

func RecDataPack(conn net.Conn) ([]byte, error) {
	recDone 		:= false
	curReadHeadSize	:= uint64(0)
	curReadBodySize	:= uint64(0)
	var buffer []byte //作为缓冲区，为了接受较长的body，可以随便设置长一点，这里也可以不限制长度，比较随意，因为后面实际读取的时候会动态扩容
	for ; !recDone ; {
		//1 获取请求报头
		readSize, err := recMsgHead(conn, curReadHeadSize, &buffer)
		if err != nil {
			curReadHeadSize += readSize
			//1.1 HEAD_ERROR代表数据获取成功，但是报头不合法，需要重新获取报头，TODO:可以考虑是不是打个日志啥的
			if err == HeadError {
				fmt.Printf("读取头部数据错误，内容为:%+v", buffer[0:curReadHeadSize])
				continue
			}
			return nil, err
		}
		curReadHeadSize += readSize

		//2 获取请求Body
		head := convertBytesToHeader(buffer)
		readSize, err = recMsgBody(conn, head, &buffer, curReadBodySize)
		if err != nil {
			curReadBodySize += readSize
			//2.1 BODY_ERROR代表数据获取成功，但是Body不合法，需要重新获取报头，TODO:可以考虑是不是打个日志啥的
			if err == BodyError {
				fmt.Printf("读取Body数据错误，内容为:%+v", buffer[curReadHeadSize:curReadBodySize])
				continue
			}
			return nil, err
		}
		curReadBodySize += readSize
		recDone = true
	}
	return buffer[HeadSize:], nil
}

/**
	函数名：recMsgHead
	功能描述：获取请求报头，这个函数可重入，如果body获取失败，可以重新获取报头
	参数1：buffer byte切片，存储从网络流中读取的数据
	参数2：curReadSize uint64，代表当前包总共读取的字节数
	返回值1：size uint64，此次读取的字节数，读取失败时为0
	返回值2：err，返回读取错误
*/
func recMsgHead(conn net.Conn, curReadSize uint64, buffer *[]byte) (size uint64, err error) {
	//1 获取字节流，直到获取到headSize大小的数据。
	headSize 	:= HeadSize
	var thisTimeReadSize uint64 = 0
	for ; curReadSize < headSize ; {
		waitReadSize	:= headSize - curReadSize
		headBuffer		:= make([]byte, waitReadSize)
		readSize, err	:= conn.Read(headBuffer)
		if err != nil {
			return thisTimeReadSize, err
		}
		if readSize <= 0 {
			continue
		}
		if uint64(len(*buffer)) <= curReadSize + uint64(readSize) {
			newBuffer := make([]byte, curReadSize + uint64(readSize))
			*buffer = append(*buffer, newBuffer...)
		}
		copy((*buffer)[curReadSize:], headBuffer[:readSize])
		curReadSize 		+= uint64(readSize)
		thisTimeReadSize 	+= uint64(readSize)
	}

	//2 校验头部是否完整
	checkErr := checkHead(*buffer)
	if !checkErr {
		//2.1 去掉头部第一个字节
		index := uint64(0)
		for ; index < curReadSize - 1 ; index++ {
			(*buffer)[index] = (*buffer)[index + 1]
		}
		//2.2 再读入一个字节
		oneByteBuffer	:= make([]byte, 1)
		for {
			readSize, err 	:= conn.Read(oneByteBuffer)
			if err != nil {
				return thisTimeReadSize, err
			}
			if readSize < 1 {
				continue
			}
			break
		}
		//2.3 把读入的字节接到之前的buffer后面
		(*buffer)[index] = oneByteBuffer[0]
		return thisTimeReadSize, HeadError
	}
	return thisTimeReadSize, nil
}

/**
* 函数名：recMsgBody
* 功能描述：获取请求Body，这个函数可重入
* 参数1：head byte切片，存储报头的内容，以便获取
* 参数2：buffer byte切片，存储从网络流中读取的数据
* 参数3：curReadSize uint64，代表当前包总共读取的字节数
* 返回值1：size uint64，此次读取的字节数，读取失败时为0
* 返回值2：err，返回读取错误
*/
func recMsgBody(conn net.Conn, head *Head, buffer *[]byte, curReadSize uint64) (size uint64, err error) {
	//1 从Head中获取Body的长度，并初始化buffer大小
	bodySize 	:= head.BodyLength

	//2 获取字节流，直到获取到bodySize大小的数据。
	var thisTimeReadSize uint64 = 0
	for ; curReadSize < bodySize ; {
		waitReadSize	:= bodySize - curReadSize
		bodyBuffer		:= make([]byte, waitReadSize)
		readSize, err	:= conn.Read(bodyBuffer)
		if err != nil {
			return thisTimeReadSize, err
		}
		if readSize <= 0 {
			continue
		}
		if uint64(len(*buffer)) <= HeadSize + curReadSize + uint64(readSize) {
			newBuffer := make([]byte, curReadSize + uint64(readSize))
			*buffer = append(*buffer, newBuffer...)
		}
		copy((*buffer)[HeadSize + curReadSize:], bodyBuffer[:readSize])
		curReadSize 		+= uint64(readSize)
		thisTimeReadSize 	+= uint64(readSize)
	}

	//3 调用实现类的checkBody方法，校验Body是否完整
	checkErr	:= checkBody(head, (*buffer)[HeadSize:HeadSize+curReadSize])
	if !checkErr {
		//3.1 去掉头部第一个字节
		index := uint64(0)
		for ; index < HeadSize + curReadSize - 1 ; index++ {
			(*buffer)[index] = (*buffer)[index + 1]
		}
		//3.2 再读入一个字节
		oneByteBuffer	:= make([]byte, 1)
		for {
			readSize, err 	:= conn.Read(oneByteBuffer)
			if err != nil {
				return thisTimeReadSize, err
			}
			if readSize < 1 {
				continue
			}
			break
		}
		//3.3 把读入的字节接到之前的buffer后面
		(*buffer)[index] = oneByteBuffer[0]
		return thisTimeReadSize, BodyError
	}
	return thisTimeReadSize, nil
}

func checkHead(buffer []byte) bool {
	head := convertBytesToHeader(buffer)
	if head.Mark != HeadMark {
		return false
	}

	prefix := buffer[0:HeadSize-16]
	md5Byte := md5.Sum(prefix)
	if !reflect.DeepEqual(md5Byte, head.HeadMd5) {
		return false
	}
	return true
}

func checkBody(head *Head, body []byte) bool {
	return reflect.DeepEqual(md5.Sum(body), head.BodyMd5)
}

func SendTcpMsg(conn net.Conn, bodyBuffer []byte) (err error) {
	head := Head{
		Mark: HeadMark,
		BodyLength: uint64(len(bodyBuffer)),
		BodyMd5: md5.Sum(bodyBuffer),
	}

	var headBuffer []byte
	headBuffer = append(headBuffer, []byte(head.Mark)...)
	headBuffer = append(headBuffer, IntToBytes(head.BodyLength)...)
	headBuffer = append(headBuffer, head.BodyMd5[0:16]...)
	md5Str := md5.Sum(headBuffer)
	headBuffer = append(headBuffer, md5Str[0:16]...)

	buffer := append(headBuffer, bodyBuffer...)
	// 这里可以考虑加锁，防止并发写。
	err		= sendMsg(conn, buffer)
	return err
}

func sendMsg (conn net.Conn, buffer []byte) (err error){
	allSize		:= 0
	nowIndex	:= 0
	bufferSize 	:= len(buffer)
	for {
		currSize, err := conn.Write(buffer[nowIndex:])
		if err != nil {
			return err
		}
		allSize += currSize
		if allSize >= bufferSize {
			break
		}
		nowIndex = allSize
	}
	return nil
}