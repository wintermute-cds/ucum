package ucum

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"UCUM_Golang/ucum"
)

func TestInterfaceImplementation_CelsiusHandler(t *testing.T){
	var _ ucum.SpecialUnitHandlerer = (*ucum.CelsiusHandler)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}
