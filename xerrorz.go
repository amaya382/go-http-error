package xerrorz

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

type HTTPErr struct {
	ErrDoc HTTPErrDoc `json:"error"`

	frame xerrors.Frame `json:"-"`
}

type HTTPErrDoc struct {
	Errors  []*InnerErr `json:"errors"`
	Code    int         `json:"code" example:"429"`
	Message string      `json:"message" example:"Rate Limit Exceeded"`

	frame xerrors.Frame `json:"-"`
}

type InnerErr struct {
	Domain       string `json:"domain" example:"usage"`                // global, {yourService}, usage,...
	Reason       string `json:"reason" example:"rateLimitExceeded"`    // invalidParameter, required,...
	Location     string `json:"location" example:""`                   // Authorization, {paramName},...
	LocationType string `json:"locationType" example:""`               // header, parameter,...
	Message      string `json:"message" example:"Rate Limit Exceeded"` // {description}
	Cause        error  `json:"-"`                                     // For error reporting

	frame xerrors.Frame `json:"-"`
}

func (e HTTPErr) Error() string {
	return fmt.Sprintf("[HTTP Status %d] %s\n", e.ErrDoc.Code, e.ErrDoc.Message)
}

func (e HTTPErr) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e HTTPErr) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.ErrDoc
}

func (e HTTPErr) Unwrap() error {
	return e.ErrDoc
}

func (e HTTPErrDoc) Error() string {
	return e.Message
}

func (e HTTPErrDoc) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e HTTPErrDoc) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.Unwrap()
}

func (e HTTPErrDoc) Unwrap() error {
	if len(e.Errors) == 0 {
		return nil
	}

	queue := []*error{}
	for _, iErr := range e.Errors[1:] {
		var err error
		err = iErr
		queue = append(queue, &err)
	}
	return ErrQueue{
		Curr:  *e.Errors[0],
		Queue: queue}
}

func (e InnerErr) Unwrap() error {
	return e.Cause
}

func (e InnerErr) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e InnerErr) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.Cause
}

func (e InnerErr) Error() string {
	return e.Message
}

func NewHTTPErr(errType ErrType, innerErrs ...*InnerErr) *HTTPErr {
	errDoc := errs[errType]
	errDoc.frame = xerrors.Caller(0)
	res := &HTTPErr{
		ErrDoc: errDoc,
		frame:  xerrors.Caller(1)}
	if innerErrs != nil {
		res.ErrDoc.Errors = innerErrs
	} else {
		res.ErrDoc.Errors = []*InnerErr{}
	}
	return res
}

func NewInnerErr(domain string, reason string, location string,
	locationType string, message string, cause error) *InnerErr {
	return &InnerErr{
		Domain:       domain,
		Reason:       reason,
		Location:     location,
		LocationType: locationType,
		Message:      message,
		Cause:        cause,
		frame:        xerrors.Caller(1)}
}

// Based on https://cloud.google.com/storage/docs/json_api/v1/status-codes#http-status-and-error-codes
type ErrType uint

const (
	// 400
	BadRequest ErrType = iota
	InvalidAltVaule
	InvalidArgument
	InvalidParameter
	ParseError
	Required
	TurnedDown

	// 401
	AuthenticationError
	NotAuthenticated
	NotAuthorized

	// 403
	AccountDisabled
	CountryBlocked
	Forbidden
	InsufficientPermissions
	SSLRequired

	// 404
	NotFound

	// 405
	MethodNotAllowed

	// 409
	Conflict

	// 410
	Gone

	// 411
	LengthRequired

	// 412
	ConditionNotMet

	// 413
	PayloadTooLarge

	// 416
	RequestedRangeNotSatisfiable

	// 429
	RateLimitExceeded
	UserRateLimitExceeded

	// 500
	InternalServerError

	// 502
	BadGateway

	// 503
	ServiceUnavailable
)

var errs = map[ErrType]HTTPErrDoc{
	BadRequest: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Bad request"},
	InvalidAltVaule: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Invalid alt value"},
	InvalidArgument: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Invalid argument"},
	InvalidParameter: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameter"},
	ParseError: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse"},
	Required: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "Required parameter or request body is missing"},
	TurnedDown: HTTPErrDoc{
		Code:    http.StatusBadRequest,
		Message: "No longer available endpoint"},
	// Tried to authenticate but authn info was not found or invalid state such as failure to parse
	AuthenticationError: HTTPErrDoc{
		Code:    http.StatusUnauthorized,
		Message: "Authentication required"}, // FIXME: or BadRequest?
	// Tried to authenticate but authn info was invalid
	NotAuthenticated: HTTPErrDoc{
		Code:    http.StatusUnauthorized,
		Message: "Authentication failed"},
	// Tried to authorize but the identified user didn't have permission to do
	NotAuthorized: HTTPErrDoc{
		Code:    http.StatusUnauthorized,
		Message: "Authorization failed"},
	AccountDisabled: HTTPErrDoc{
		Code:    http.StatusForbidden,
		Message: "Account has been disabled"},
	CountryBlocked: HTTPErrDoc{
		Code:    http.StatusForbidden,
		Message: "Restricted by law with your country"},
	Forbidden: HTTPErrDoc{
		Code:    http.StatusForbidden,
		Message: "Not allowed endpoint"},
	InsufficientPermissions: HTTPErrDoc{
		Code:    http.StatusForbidden,
		Message: "Insufficient permissions"},
	SSLRequired: HTTPErrDoc{
		Code:    http.StatusForbidden,
		Message: "SSL is required"},
	NotFound: HTTPErrDoc{
		Code:    http.StatusNotFound,
		Message: "Not found"},
	MethodNotAllowed: HTTPErrDoc{
		Code:    http.StatusMethodNotAllowed,
		Message: "Not allowed method"},
	Conflict: HTTPErrDoc{
		Code:    http.StatusConflict,
		Message: "Conflict"},
	Gone: HTTPErrDoc{
		Code:    http.StatusGone,
		Message: "Resources or session has gone"},
	LengthRequired: HTTPErrDoc{
		Code:    http.StatusLengthRequired,
		Message: "Content-Length header is required"},
	ConditionNotMet: HTTPErrDoc{
		Code:    http.StatusPreconditionFailed,
		Message: "Pre-condition did not hold"},
	PayloadTooLarge: HTTPErrDoc{
		Code:    http.StatusRequestEntityTooLarge,
		Message: "Too large payload"},
	RequestedRangeNotSatisfiable: HTTPErrDoc{
		Code:    http.StatusRequestedRangeNotSatisfiable,
		Message: "Requested range cannot be satisfied"},
	RateLimitExceeded: HTTPErrDoc{
		Code:    http.StatusTooManyRequests,
		Message: "Rate quota was exceeded"},
	UserRateLimitExceeded: HTTPErrDoc{
		Code:    http.StatusTooManyRequests,
		Message: "Per-user rate quota was exceeded"},
	InternalServerError: HTTPErrDoc{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error"},
	BadGateway: HTTPErrDoc{
		Code:    http.StatusBadGateway,
		Message: "Bad gateway"},
	ServiceUnavailable: HTTPErrDoc{
		Code:    http.StatusServiceUnavailable,
		Message: "Temporarily service unavailable"},
}
