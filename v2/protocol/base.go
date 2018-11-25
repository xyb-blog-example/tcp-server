package protocol
//base.go文件，定义了一些协议需要用的结构体，类似C/C++的头文件

type Head struct {
	Mark			string
	BodyLength		uint64
	BodyMd5			[16]byte
	HeadMd5			[16]byte
}

const HeadSize = uint64(44)

const HeadMark = "HAHA"