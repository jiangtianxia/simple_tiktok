package utils

import (
	"fmt"

	"github.com/speps/go-hashids"
)

var h *hashids.HashID

func InitHashID(salt string, minLength int) error {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength

	hs, err := hashids.NewWithData(hd)
	if err != nil {
		return err
	}
	h = hs
	return nil
}

// EncodeID 将 uint 类型的 ID 加密成字符串
func EncodeID(id uint) (string, error) {
	encodedID, err := h.Encode([]int{int(id)})
	if err != nil {
		return "", err
	}
	return encodedID, nil
}

// DecodeID 将加密后的字符串解密成 uint 类型的 ID
func DecodeID(encodedID string) (uint, error) {
	decodedID, err := h.DecodeWithError(encodedID)
	if err != nil {
		return 0, err
	}
	if len(decodedID) == 0 {
		return 0, fmt.Errorf("invalid encoded ID")
	}
	return uint(decodedID[0]), nil
}
