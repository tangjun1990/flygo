package xhashid

import (
	"log"
	"testing"
)

func TestHashIds_Encode(t *testing.T) {
	hash, _ := New("mysaltxxxx", 10)
	e, _ := hash.Encode(123456)
	if e != "AJk77eaEkB" {
		log.Fatal("not match")
	}
}

func TestHashIds_Decode(t *testing.T) {
	hash, _ := New("mysaltxxxx", 10)
	d, _ := hash.Decode("AJk77eaEkB")
	if d != 123456 {
		log.Fatal("not match")
	}
}
