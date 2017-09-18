package ucum

import (
	"UCUM_Golang/ucum"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInterfaceImplementation_CelsiusHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.CelsiusHandler)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}

func TestInterfaceImplementation_FahrenheitHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.FahrenheitHandler)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}

func TestInterfaceImplementation_HoldingHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.HoldingHandler)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}
