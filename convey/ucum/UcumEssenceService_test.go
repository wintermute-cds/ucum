package ucum

import (
	"testing"
	"os"
	"UCUM_Golang/ucum"
	"github.com/smartystreets/goconvey/convey"
)

var definitions string
var tests string

func TestService(t *testing.T){
	convey.Convey("TestStringAsIntegerDecimal", t, func() {
		definitions := os.Getenv("GOPATH") + "/src/UCUM_Golang/ucum/terminology_data/ucum-essence.xml"
		service, err := ucum.GetInstanceOfOpenEhrTerminologyService(definitions)
		convey.So(err, convey.ShouldBeNil)
		convey.So(service, convey.ShouldNotBeNil)
	})
}

type
