package xdecimal

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMulWithString(t *testing.T) {
	Convey("浮点数乘法测试：", t, func() {
		s, err := MulWithString("0.57", "100")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "57")
	})
}

func TestAddWithString(t *testing.T) {
	Convey("浮点数加法测试：", t, func() {
		s, err := AddWithString("57.01", "45.121")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "102.131")
	})
}

func TestDivWithString(t *testing.T) {
	Convey("浮点数除法测试：", t, func() {
		s, err := DivWithString("57", "100")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "0.57")
	})
}

func TestSubWithString(t *testing.T) {
	Convey("浮点数减法测试：", t, func() {
		s, err := SubWithString("0.57", "0.00001")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "0.56999")
	})
}
