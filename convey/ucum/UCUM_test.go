package ucum

import (
	"testing"
	"UCUM_Golang/ucum"
	"github.com/smartystreets/goconvey/convey"
)

func TestInterfaceImplementation_UCUM(t *testing.T){
	var _ ucum.Concepter = (*ucum.Unit)(nil)
	var _ ucum.Concepter = (*ucum.Prefix)(nil)
	var _ ucum.Concepter = (*ucum.BaseUnit)(nil)
	var _ ucum.Concepter = (*ucum.DefinedUnit)(nil)
	var _ ucum.Uniter = (*ucum.BaseUnit)(nil)
	var _ ucum.Uniter = (*ucum.DefinedUnit)(nil)
	convey.Convey("Formal inheritance test", t, func() {
		convey.So(true, convey.ShouldBeTrue)
	})
}
