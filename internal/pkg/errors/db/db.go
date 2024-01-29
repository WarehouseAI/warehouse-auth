package errors

type DBErrorType string

const (
	DbExist      DBErrorType = "exist"
	DbNotFound   DBErrorType = "not_found"
	DbSystem     DBErrorType = "system"
	DbForeignKey DBErrorType = "foreign_key"
)

type Error interface {
	error
	GetStatus() DBErrorType
	GetTrace() string
}

type DBError struct {
	Status DBErrorType
	Trace  string
	Err    error
}

func (dep DBError) Error() string {
	return dep.Err.Error()
}

func (dep DBError) GetStatus() DBErrorType {
	return dep.Status
}

func (dep DBError) GetTrace() string {
	return dep.Trace
}

func NewDBError(errType DBErrorType, trace string, err error) error {
	return DBError{errType, trace, err}
}
