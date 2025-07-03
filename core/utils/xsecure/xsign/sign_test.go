package xsign

import (
	"log"
	"testing"
)

func TestGenSign(t *testing.T) {
	var param = map[string]string{
		"id":   "1",
		"name": "tangjun",
	}
	sign := GenSign("secret1212", param)
	if sign != "549f1d1efd864558357922a0a708d908" {
		log.Fatal("not match")
	}
}
