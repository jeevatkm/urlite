package tpl

import (
	"bytes"
	"errors"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var FuncMap template.FuncMap

func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

func IsSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

func SafeCSS(text string) template.CSS {
	return template.CSS(text)
}

func SafeURL(text string) template.URL {
	return template.URL(text)
}

func FriendlyDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	if t.Year() == time.Now().Year() {
		return t.Format("Jan 2, 3:04:05 pm")
	}
	return t.Format("Jan 2, 2006, 3:04:05 pm")
}

func ToCommaStr(v []string) string {
	return strings.Join(v, ", ")
}

func Add(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	if av.Kind() != bv.Kind() {
		return nil, errors.New("Different kinds, can't add them.")
	}

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() + bv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return av.Uint() + bv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return av.Float() + bv.Float(), nil
	// case reflect.String:
	//     return av.String() + bv.String(), nil
	default:
		return nil, errors.New("Type does not support addition.")
	}
}

func NumToStr(n int64, sep rune) string {
	s := strconv.FormatInt(n, 10)

	startOffset := 0
	var buff bytes.Buffer
	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}

	l := len(s)
	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}

	for i := startOffset; i < l; i++ {
		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++

		buff.WriteByte(s[i])
	}

	return buff.String()
}

func init() {
	FuncMap = template.FuncMap{
		"safeHTML":         SafeHTML,
		"safeCSS":          SafeCSS,
		"safeURL":          SafeURL,
		"friendlyDateTime": FriendlyDateTime,
		"toCommaStr":       ToCommaStr,
		"add":              Add,
		"isSet":            IsSet,
		"num2str":          NumToStr,
	}
}
