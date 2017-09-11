package ucum

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"UCUM_Golang/ucum"
	"strconv"
	"github.com/smartystreets/assertions/should"
)

const MAX_INT = 0x7fffffff
const MIN_INT = 0x80000000

func TestRoundtripIntegerDecimal(t *testing.T){
	convey.Convey("TestRoundtripIntegerDecimal", t, func(){
		d,err := ucum.NewDecimal(strconv.Itoa(0))
		convey.So(err, should.BeNil)
		i,err := d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, 0)
		d,err = ucum.NewDecimal(strconv.Itoa(1))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, 1)
		d,err = ucum.NewDecimal(strconv.Itoa(2))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, 2)
		d,err = ucum.NewDecimal(strconv.Itoa(64))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, 64)
		d,err = ucum.NewDecimal(strconv.Itoa(MAX_INT))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, MAX_INT)
		d,err = ucum.NewDecimal(strconv.Itoa(-1))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, -1)
		d,err = ucum.NewDecimal(strconv.Itoa(-2))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, -2)
		d,err = ucum.NewDecimal(strconv.Itoa(-64))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, -64)
		d,err = ucum.NewDecimal(strconv.Itoa(MIN_INT))
		convey.So(err, should.BeNil)
		i,err = d.AsInteger()
		convey.So(err, should.BeNil)
		convey.So(i, should.Equal, MIN_INT)
	})
}
