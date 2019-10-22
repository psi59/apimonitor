package rshttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	return strings.ToUpper(string(method))
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

func (method *Method) UnmarshalJSON(data []byte) error {
	var str string
	if err := jsoniter.Unmarshal(data, &str); err != nil {
		return errors.WithStack(err)
	}
	*method = Method(strings.ToUpper(str))
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

func (contentType ContentType) GetBodyFromResponse(res *http.Response) (interface{}, error) {
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var body interface{}
	switch contentType {
	case MIMEApplicationJSON:
		body = make(map[string]interface{})
	case MIMETextHTML:
		body = ""
	default:
		body = make(map[string]interface{})
	}

	if err := jsoniter.Unmarshal(rawBody, &body); err != nil {
		return nil, errors.WithStack(err)
	}

	return body, nil
}
