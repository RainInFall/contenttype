package contenttype

import (
	"regexp"
	"testing"

	"github.com/RainInFall/assert"
)

func TestFormat(t *testing.T) {
	assert.Init(t)

	func() {
		str, err := Format(&ContentType{Type: "text/html"})
		assert.Ok(nil == err)
		assert.Ok("text/html" == str)
	}()

	func() {
		str, err := Format(&ContentType{Type: "image/svg+xml"})
		assert.Ok(nil == err)
		assert.Ok("image/svg+xml" == str)
	}()

	func() {
		str, err := Format(&ContentType{
			Type: "text/html",
			Parameters: map[string]string{
				"charset": "utf-8",
			},
		})
		assert.Ok(nil == err)
		assert.Ok("text/html; charset=utf-8" == str)
	}()

	func() {
		str, err := Format(&ContentType{
			Type: "text/html",
			Parameters: map[string]string{
				"foo": "bar or \"baz\"",
			},
		})
		assert.Ok(nil == err)
		assert.Ok("text/html; foo=\"bar or \\\"baz\\\"\"" == str)
	}()

	func() {
		str, err := Format(&ContentType{
			Type: "text/html",
			Parameters: map[string]string{
				"foo": "",
			},
		})
		assert.Ok(nil == err)
		assert.Ok("text/html; foo=\"\"" == str)
	}()

	func() {
		str, err := Format(&ContentType{
			Type: "text/html",
			Parameters: map[string]string{
				"charset": "utf-8",
				"foo":     "bar",
				"bar":     "baz",
			},
		})
		assert.Ok(nil == err)
		assert.Ok("text/html; bar=baz; charset=utf-8; foo=bar" == str)
	}()

	func() {
		_, err := Format(nil)
		assert.Ok(nil != err && regexp.MustCompile("argument ContentType is required").MatchString(err.Error()))
	}()

	func() {
		_, err := Format(&ContentType{})
		assert.Ok(nil != err && regexp.MustCompile("invalid type").MatchString(err.Error()))
	}()

	func() {
		_, err := Format(&ContentType{Type: "text/"})
		assert.Ok(nil != err && regexp.MustCompile("invalid type").MatchString(err.Error()))
	}()

	func() {
		_, err := Format(&ContentType{Type: " text/html"})
		assert.Ok(nil != err && regexp.MustCompile("invalid type").MatchString(err.Error()))
	}()

	func() {
		_, err := Format(&ContentType{
			Type: "image/svg",
			Parameters: map[string]string{
				"foo/": "bar",
			},
		})
		assert.Ok(nil != err && regexp.MustCompile("invalid parameter name").MatchString(err.Error()))
	}()

	func() {
		_, err := Format(&ContentType{
			Type: "image/svg",
			Parameters: map[string]string{
				"foo": "bar\u0000",
			},
		})
		assert.Ok(nil != err && regexp.MustCompile("invalid parameter value").MatchString(err.Error()))
	}()
}
