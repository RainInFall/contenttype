package contenttype

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

//go:generate js-like object string string

type ContentType struct {
	typ        string
	parameters Objectstringstring
}

func Parameters(a map[string]string) Objectstringstring {
	return a
}

var paramRegExp = regexp.MustCompile("; *([!#$%&'\\*\\+\\-\\.\\^_`\\|~0-9A-Za-z]+) *= *(\"(?:[\\x{000b}\\x{0020}\\x{0021}\\x{0023}-\\x{005b}\\x{005d}-\\x{007e}\\x{0080}-\\x{00ff}]|\\\\[\\x{000b}\\x{0020}-\\x{00ff}])*\"|[!#$%&'\\*\\+\\-\\.\\^_`\\|~0-9A-Za-z]+) *")
var textRegExp = regexp.MustCompile("^[\\x{000b}\\x{0020}-\\x{007e}\\x{0080}-\\x{00ff}]+$")
var tokenRegExp = regexp.MustCompile("^[!#$%&'\\*\\+\\-\\.\\^_`\\|~0-9A-Za-z]+$")
var qescRegExp = regexp.MustCompile("\\\\([\\x{000b}\\x{0020}-\\x{00ff}])")
var quoteRegExp = regexp.MustCompile("([\\\\\"])")
var typeRegExp = regexp.MustCompile("^[!#$%&'\\*\\+\\-\\.\\^_`\\|~0-9A-Za-z]+\\/[!#$%&'\\*\\+\\-\\.\\^_`\\|~0-9A-Za-z]+$")

/*
Format object to media type
*/
func Format(contentType *ContentType) (string, error) {
	if !typeRegExp.MatchString(contentType.typ) {
		return "", errors.New("invalid type")
	}

	result := contentType.typ

	for _, param := range contentType.parameters.Keys().Sort() {
		if !tokenRegExp.MatchString(param) {
			return "", errors.New("invalid paremeter name")
		}
		if value, err := qstring(contentType.parameters[param]); nil == err {
			result += "; " + param + "=" + value
		} else {
			return "", err
		}
	}

	return result, nil
}

/*
Parse
*/

func ParseRequest(req interface{}) (*ContentType, error) {
	headerField := reflect.ValueOf(req).FieldByName("Header")
	if !headerField.IsValid() || !headerField.Type().ConvertibleTo(reflect.TypeOf((http.Header)(nil))) {
		return nil, errors.New("content-type header is missing from object")
	}
	return ParseHeader(headerField.Interface().(http.Header))
}

/*
ParseHeader parse media type from http header
*/
func ParseHeader(header http.Header) (*ContentType, error) {
	str := header.Get("content-type")
	if "" == str {
		return nil, errors.New("content-type header is missing from object")
	}
	return Parse(str)
}

/*
Parse media type to object
*/
func Parse(str string) (*ContentType, error) {
	index := strings.Index(str, ";")
	var typ string
	if -1 != index {
		typ = strings.TrimSpace(str[0:index])
	} else {
		typ = strings.TrimSpace(str)
	}

	if !typeRegExp.MatchString(typ) {
		return nil, errors.New("invalid media type")
	}

	typ = strings.ToLower(typ)
	parameters := make(map[string]string)

	for _, match := range paramRegExp.FindAllStringSubmatchIndex(str, -1) {
		if match[0] != index {
			return nil, errors.New("invalid parameter format")
		}

		index = match[1]
		key := strings.ToLower(str[match[2]:match[3]])
		value := str[match[4]:match[5]]

		if '"' == value[0] {
			value = qescRegExp.ReplaceAllStringFunc(
				value[1:len(value)-1],
				func(match string) string {
					return qescRegExp.FindStringSubmatch(match)[1]
				})
		}

		parameters[key] = value
	}

	if -1 != index && len(str) != index {
		return nil, errors.New("invalid parameter format")
	}

	return &ContentType{
		typ,
		parameters,
	}, nil
}

func qstring(str string) (string, error) {
	if tokenRegExp.MatchString(str) {
		return str, nil
	}

	if len(str) > 0 && !textRegExp.MatchString(str) {
		return str, errors.New("invalid parameter value")
	}

	return "\"" + quoteRegExp.ReplaceAllStringFunc(str, func(match string) string {
		return "\\" + quoteRegExp.FindStringSubmatch(match)[1]
	}) + "\"", nil
}
