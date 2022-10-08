package ft

import (
	"bytes"
)

type pkcs5 uint

func newPkcs5(size int) pkcs5 {
	return pkcs5(size)
}

func (p pkcs5) BlockSize() int {
	return int(p)
}

func (p pkcs5) Pad(src []byte) (dst []byte) {
	count := p.BlockSize() - len(src)%p.BlockSize()
	padded := bytes.Repeat([]byte{byte(count)}, count)
	dst = append(src, padded...)

	return
}

func (p pkcs5) Unpad(src []byte) (dst []byte, err error) {
	length := len(src)
	// 通过最后一位找到填充了几，填充了几个
	padding := int(src[length-1])
	dst = src[:(length - padding)]

	return
}
