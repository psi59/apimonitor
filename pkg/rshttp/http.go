package rshttp

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/realsangil/apimonitor/pkg/rserrors"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"

	MIMEApplicationJSON       = "application/json"
	MIMEApplicationJavaScript = "application/javascript"
	MIMEApplicationXML        = "application/xml"
	MIMETextXML               = "text/xml"
	MIMEApplicationForm       = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf   = "application/protobuf"
	MIMEApplicationMsgpack    = "application/msgpack"
	MIMETextHTML              = "text/html"
	MIMETextPlain             = "text/plain"
	MIMEMultipartForm         = "multipart/form-data"
	MIMEOctetStream           = "application/octet-stream"
)

var (
	MIMEApplicationJSONCharsetUTF8       = appendUTF8IntoMIME(MIMEApplicationJSON)
	MIMEApplicationJavaScriptCharsetUTF8 = appendUTF8IntoMIME(MIMEApplicationJavaScript)
	MIMEApplicationXMLCharsetUTF8        = appendUTF8IntoMIME(MIMEApplicationXML)
	MIMETextXMLCharsetUTF8               = appendUTF8IntoMIME(MIMETextXML)
	MIMETextHTMLCharsetUTF8              = appendUTF8IntoMIME(MIMETextHTML)
	MIMETextPlainCharsetUTF8             = appendUTF8IntoMIME(MIMETextPlain)

	appendUTF8IntoMIME = func(mime string) string {
		return fmt.Sprintf("%s; charset=UTF-8", mime)
	}
)

var regexpEndpointPath = regexp.MustCompile(`^[\/a-zA-Z0-9_.\-~!$&'()*+,;=:@]+$`)

type EndpointPath string

func (endpointPath EndpointPath) String() string {
	return string(endpointPath)
}

func (endpointPath EndpointPath) Validate() error {
	if regexpEndpointPath.MatchString(endpointPath.String()) {
		return nil
	}
	return errors.Wrap(rserrors.ErrInvalidParameter, "EndpointPath")
}

type Method string

func (method Method) String() string {
	return string(method)
}

func (method Method) Validate() error {
	switch method {
	case MethodGet:
	case MethodHead:
	case MethodPost:
	case MethodPut:
	case MethodPatch:
	case MethodDelete:
	case MethodConnect:
	case MethodOptions:
	case MethodTrace:
	default:
		return errors.Wrap(rserrors.ErrInvalidParameter, "Method")
	}
	return nil
}

var regexpContentType = regexp.MustCompile(`(text|application|multipart)/(javascript|json|x-www-form-urlencoded|octet-stream|form-data|xml)(;(.+))?`)

type ContentType string

func (contentType ContentType) String() string {
	return string(contentType)
}

func (contentType ContentType) Validate() error {
	if regexpContentType.MatchString(contentType.String()) {
		return nil
	}
	return errors.Wrap(rserrors.ErrInvalidParameter, "ContentType")
}
