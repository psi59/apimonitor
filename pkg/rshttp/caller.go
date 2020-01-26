package rshttp

import (
	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

func Do(request *Request) (*Response, error) {
	if rsvalid.IsZero(request) {
		return nil, errors.WithStack(rserrors.ErrInvalidParameter)
	}
	return request.execute()
}
