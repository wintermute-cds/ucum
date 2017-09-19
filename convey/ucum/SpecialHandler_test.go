package ucum

import (
	"ucum"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInterfaceImplementation_CelsiusHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.CelsiusHandler)(nil)
	Convey("Formal inheritance test", t, func() {
		So(true, ShouldBeTrue)
	})
}

func TestInterfaceImplementation_FahrenheitHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.FahrenheitHandler)(nil)
	Convey("Formal inheritance test", t, func() {
		So(true, ShouldBeTrue)
	})
}

func TestInterfaceImplementation_HoldingHandler(t *testing.T) {
	var _ ucum.SpecialUnitHandlerer = (*ucum.HoldingHandler)(nil)
	Convey("Formal inheritance test", t, func() {
		So(true, ShouldBeTrue)
	})
}
