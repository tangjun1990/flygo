package xdecimal

import (
	"github.com/shopspring/decimal"
)

// MulWithString returns s1 * s2.
func MulWithString(s1, s2 string) (s3 string, err error) {
	d1, err := decimal.NewFromString(s1)
	if err != nil {
		return
	}
	d2, err := decimal.NewFromString(s2)
	if err != nil {
		return
	}
	s3 = d1.Mul(d2).String()
	return
}

// DivWithString returns s1 / s2.
func DivWithString(s1, s2 string) (s3 string, err error) {
	d1, err := decimal.NewFromString(s1)
	if err != nil {
		return
	}
	d2, err := decimal.NewFromString(s2)
	if err != nil {
		return
	}
	s3 = d1.Div(d2).String()
	return
}

// AddWithString returns s1 + s2.
func AddWithString(s1, s2 string) (s3 string, err error) {
	d1, err := decimal.NewFromString(s1)
	if err != nil {
		return
	}
	d2, err := decimal.NewFromString(s2)
	if err != nil {
		return
	}
	s3 = d1.Add(d2).String()
	return
}

// SubWithString returns s1 - s2.
func SubWithString(s1, s2 string) (s3 string, err error) {
	d1, err := decimal.NewFromString(s1)
	if err != nil {
		return
	}
	d2, err := decimal.NewFromString(s2)
	if err != nil {
		return
	}
	s3 = d1.Sub(d2).String()
	return
}
