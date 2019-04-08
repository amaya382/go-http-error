package httperror

type HTTPErr struct {
	ErrDoc HTTPErrDoc `json:"error"`
}

type HTTPErrDoc struct {
	Errors  []*InnerErr `json:"errors"`
	Code    int         `json:"code" example:"429"`
	Message string      `json:"message" example:"Rate Limit Exceeded"`
}

type InnerErr struct {
	Domain       string `json:"domain" example:"usage"`                // global, {yourService}, usage,...
	Reason       string `json:"reason" example:"rateLimitExceeded"`    // invalidParameter, required,...
	Location     string `json:"location" example:""`                   // Authorization, {paramName},...
	LocationType string `json:"locationType" example:""`               // header, parameter,...
	Message      string `json:"message" example:"Rate Limit Exceeded"` // {description}
}

func (e HTTPErr) Error() string {
	return e.ErrDoc.Message
}

func (e HTTPErrDoc) Error() string {
	return e.Message
}

func (e InnerErr) Error() string {
	return e.Message
}

func NewHTTPErr(errType ErrType, innerErrs ...*InnerErr) *HTTPErr {
	res := &HTTPErr{
		ErrDoc: errs[errType]}
	if innerErrs != nil {
		res.ErrDoc.Errors = innerErrs
	} else {
		res.ErrDoc.Errors = []*InnerErr{}
	}
	return res
}

func NewInnerErr(domain string, reason string, location string, locationType string, message string) *InnerErr {
	return &InnerErr{
		Domain:       domain,
		Reason:       reason,
		Location:     location,
		LocationType: locationType,
		Message:      message}
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
		Code:    400,
		Message: "Bad request"},
	InvalidAltVaule: HTTPErrDoc{
		Code:    400,
		Message: "Invalid alt value"},
	InvalidArgument: HTTPErrDoc{
		Code:    400,
		Message: "Invalid argument"},
	InvalidParameter: HTTPErrDoc{
		Code:    400,
		Message: "Invalid parameter"},
	ParseError: HTTPErrDoc{
		Code:    400,
		Message: "Failed to parse"},
	Required: HTTPErrDoc{
		Code:    400,
		Message: "Required parameter or request body is missing"},
	TurnedDown: HTTPErrDoc{
		Code:    400,
		Message: "No longer available endpoint"},
	AuthenticationError: HTTPErrDoc{
		Code:    401,
		Message: "Authentication failed"},
	NotAuthenticated: HTTPErrDoc{
		Code:    401,
		Message: "Authentication required"},
	AccountDisabled: HTTPErrDoc{
		Code:    403,
		Message: "Account has been disabled"},
	CountryBlocked: HTTPErrDoc{
		Code:    403,
		Message: "Restricted by law with your country"},
	Forbidden: HTTPErrDoc{
		Code:    403,
		Message: "Not allowed endpoint"},
	InsufficientPermissions: HTTPErrDoc{
		Code:    403,
		Message: "Insufficient permissions"},
	SSLRequired: HTTPErrDoc{
		Code:    403,
		Message: "SSL is required"},
	NotFound: HTTPErrDoc{
		Code:    404,
		Message: "Not found"},
	MethodNotAllowed: HTTPErrDoc{
		Code:    405,
		Message: "Not allowed method"},
	Conflict: HTTPErrDoc{
		Code:    409,
		Message: "Conflict"},
	Gone: HTTPErrDoc{
		Code:    410,
		Message: "Resources or session has gone"},
	LengthRequired: HTTPErrDoc{
		Code:    411,
		Message: "Content-Length header is required"},
	ConditionNotMet: HTTPErrDoc{
		Code:    412,
		Message: "Pre-condition did not hold"},
	PayloadTooLarge: HTTPErrDoc{
		Code:    413,
		Message: "Too large payload"},
	RequestedRangeNotSatisfiable: HTTPErrDoc{
		Code:    416,
		Message: "Requested range cannot be satisfied"},
	RateLimitExceeded: HTTPErrDoc{
		Code:    429,
		Message: "Rate quota was exceeded"},
	UserRateLimitExceeded: HTTPErrDoc{
		Code:    429,
		Message: "Per-user rate quota was exceeded"},
	InternalServerError: HTTPErrDoc{
		Code:    500,
		Message: "Internal server error"},
	BadGateway: HTTPErrDoc{
		Code:    502,
		Message: "Bad gateway"},
	ServiceUnavailable: HTTPErrDoc{
		Code:    503,
		Message: "Temporarily service unavailable"},
}
