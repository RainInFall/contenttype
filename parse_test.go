package contenttype

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/RainInFall/assert"
)

var invalidTypes = []string{
	" ",
	"null",
	"undefined",
	"/",
	"text / plain",
	"text/;plain",
	"text/\"plain",
	"text/p£ain",
	"text/(plain)",
	"text/@plain",
	"text/plain,wrong",
}

func TestParse(t *testing.T) {
	assert.Init(t)

	func() {
		typ, err := Parse("text/html")
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
	}()

	func() {
		typ, err := Parse("image/svg+xml")
		assert.Ok(nil == err)
		assert.Ok("image/svg+xml" == typ.typ)
	}()

	func() {
		typ, err := Parse(" text/html")
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
	}()

	func() {
		typ, err := Parse("text/html; charset=utf-8; foo=bar")
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"charset": "utf-8",
			"foo":     "bar",
		})))
	}()

	func() {
		typ, err := Parse("text/html ; charset=utf-8 ; foo=bar")
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"charset": "utf-8",
			"foo":     "bar",
		})))
	}()

	func() {
		typ, err := Parse("IMAGE/SVG+XML")
		assert.Ok(nil == err)
		assert.Ok("image/svg+xml" == typ.typ)
	}()

	func() {
		typ, err := Parse("text/html; Charset=UTF-8")
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"charset": "UTF-8",
		})))
	}()

	func() {
		typ, err := Parse("text/html; charset=\"UTF-8\"")
		assert.Ok(err == nil)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"charset": "UTF-8",
		})))
	}()

	func() {
		typ, err := Parse("text/html; charset = \"UT\\F-\\\\\\\"8\\\"\"")
		assert.Ok(nil == err)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"charset": "UTF-\\\"8\"",
		})))
	}()

	func() {
		typ, err := Parse("text/html; param=\"charset=\\\"utf-8\\\"; foo=bar\"; bar=foo")
		assert.Ok(nil == err)
		assert.Ok(typ.parameters.Equals(Parameters(map[string]string{
			"param": "charset=\"utf-8\"; foo=bar",
			"bar":   "foo",
		})))
	}()

	for _, v := range invalidTypes {
		_, err := Parse(v)
		assert.Ok(nil != err && regexp.MustCompile("invalid media type").MatchString(err.Error()))
	}

	func() {
		_, err := Parse("text/plain; foo=\"bar")
		assert.Ok(nil != err && regexp.MustCompile("invalid parameter format").MatchString(err.Error()))
	}()
	func() {
		_, err := Parse("text/plain; profile=http://localhost; foo=bar")
		assert.Ok(nil != err && regexp.MustCompile("invalid parameter format").MatchString(err.Error()))
	}()
	func() {
		_, err := Parse("text/plain; profile=http://localhost")
		assert.Ok(nil != err && regexp.MustCompile("invalid parameter format").MatchString(err.Error()))
	}()

	func() {
		header := make(http.Header)
		header.Set("content-type", "text/html")
		typ, err := ParseHeader(header)
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
	}()

	func() {
		req := http.Request{
			Header: http.Header{},
		}
		_, err := ParseRequest(req)
		assert.Ok(nil != err && regexp.MustCompile("content-type header is missing").MatchString(err.Error()))
	}()

	func() {
		_, err := ParseRequest(nil)
		assert.Ok(nil != err && regexp.MustCompile("content-type header is missing").MatchString(err.Error()))
	}()

	func() {
		req := http.Request{
			Header: http.Header{
				"Content-Type": []string{"text/html"},
			},
		}
		typ, err := ParseRequest(req)
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
	}()

	func() {
		res := TempResponse{http.Header{
			"Content-Type": []string{"text/html"},
		}}

		typ, err := ParserResponse(res)
		assert.Ok(nil == err)
		assert.Ok("text/html" == typ.typ)
	}()

	func() {
		res := TempResponse{}

		_, err := ParserResponse(res)
		assert.Ok(nil != err && regexp.MustCompile("content-type header is missing").MatchString(err.Error()))
	}()
}

type TempResponse struct {
	header http.Header
}

func (res TempResponse) Header() http.Header {
	return res.header
}
