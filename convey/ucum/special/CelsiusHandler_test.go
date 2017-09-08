package special

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"UCUM_Golang/ucum/special"
)

func TestInterfaceImplementation_CelsiusHandler(t *testing.T){
	var _ special.SpecialUnitHandlerer = (*special.CelsiusHandler)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}
