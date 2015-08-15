package middleware

import (
	"bytes"
	"log"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var minifier *minify.Minify
var mt []*regexp.Regexp // Media Type

func init() {
	minifier = minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	mt = []*regexp.Regexp{}
	mt = append(mt, regexp.MustCompile("[text|application]/html"))
	mt = append(mt, regexp.MustCompile("text/[css|stylesheet]"))
	mt = append(mt, regexp.MustCompile("[text|application]/javascript"))
}

type minifyWriter struct {
	http.ResponseWriter
	Code        int
	Body        *bytes.Buffer
	wroteHeader bool
}

func (m *minifyWriter) Header() http.Header {
	return m.ResponseWriter.Header()
}

func (m *minifyWriter) WriteHeader(code int) {
	if !m.wroteHeader {
		m.Code = code
		m.wroteHeader = true
		m.ResponseWriter.WriteHeader(code)
	}
}

func (m *minifyWriter) Write(b []byte) (int, error) {
	if !m.wroteHeader {
		m.WriteHeader(http.StatusOK)
	}
	if m.Body != nil {
		m.Body.Write(b)
	}
	return len(b), nil
}

func isValidMediaType(t string) bool {
	for _, v := range mt {
		if v.MatchString(t) {
			return true
		}
	}

	return false
}

func MinifyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &minifyWriter{
			ResponseWriter: w,
			Body:           new(bytes.Buffer),
		}

		h.ServeHTTP(mw, r)

		ct := w.Header().Get("Content-Type")
		if isValidMediaType(ct) {
			log.Println("Inside minify")
			rb, err := minify.Bytes(minifier, ct, mw.Body.Bytes())
			if err != nil {
				_ = err // unsupported mediatype error or internal
			}

			w.Write(rb)
		} else {
			w.Write(mw.Body.Bytes())
		}
	})
}
