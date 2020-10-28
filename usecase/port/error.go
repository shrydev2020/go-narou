package port

type ApplicationError interface {
	Error() string
	Internal() bool
}

type portError struct {
	Message string
	Code    ErrorCode
}

func (u portError) Error() string {
	return u.Message
}

func (u portError) ErrorCode() ErrorCode {
	return u.Code
}

func (u portError) Internal() bool {
	return true
}

func NewPortError(e error, code ErrorCode) ApplicationError {
	return &portError{
		Message: e.Error(),
		Code:    code,
	}
}

type ErrorCode string

func (ec ErrorCode) String() string {
	return string(ec)
}

const (
	NovelNotFoundAtSite    ErrorCode = "0"
	NovelDoesntExistInPath ErrorCode = "1"
	RepositoryError        ErrorCode = "3"
	CrawlerError           ErrorCode = "4"
	EpubError              ErrorCode = "5"
	InvalidParam           ErrorCode = "900"
	UnHandledError         ErrorCode = "999" // Error codes for internal error
)
