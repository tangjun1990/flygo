package xhashid

import (
	"errors"
	"github.com/speps/go-hashids/v2"
)

type hashIds struct {
	hashID *hashids.HashID
}

func New(salt string, minLength int) (*hashIds, error) {
	hashID := hashids.NewData()
	hashID.Salt = salt
	hashID.MinLength = minLength

	h, err := hashids.NewWithData(hashID)

	return &hashIds{
		hashID: h,
	}, err
}

func (h *hashIds) Encode(num int) (string, error) {
	return h.hashID.Encode([]int{num})
}

func (h *hashIds) Decode(hash string) (int, error) {
	d, err := h.hashID.DecodeWithError(hash)
	if err != nil {
		return 0, err
	}
	if len(d) == 0 {
		return 0, errors.New("decode error with empty slice")
	}
	return d[0], nil
}
