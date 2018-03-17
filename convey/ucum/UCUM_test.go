package ucum

import (
	"github.com/bertverhees/ucum"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInterfaceImplementation_UCUM(t *testing.T) {
	var _ ucum.Concepter = (*ucum.Unit)(nil)
	var _ ucum.Concepter = (*ucum.Prefix)(nil)
	var _ ucum.Concepter = (*ucum.BaseUnit)(nil)
	var _ ucum.Concepter = (*ucum.DefinedUnit)(nil)
	var _ ucum.Uniter = (*ucum.BaseUnit)(nil)
	var _ ucum.Uniter = (*ucum.DefinedUnit)(nil)
	Convey("Formal inheritance test", t, func() {
		So(true, ShouldBeTrue)
	})
}
