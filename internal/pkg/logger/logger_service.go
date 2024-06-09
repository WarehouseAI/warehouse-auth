package logger

import (
	"fmt"

	"github.com/warehouse/auth-service/internal/pkg/errors"
	"github.com/warehouse/auth-service/internal/pkg/errors/repository_errors"

	"go.uber.org/zap"
)

type (
	Logger interface {
		Named(name string) Logger

		Zap() *zap.Logger

		Error(err, reason error) error
		ErrorRepo(err, reason error, query string) error
		ServiceErrorWithReason(err *errors.Error, reason string) *errors.Error
		ServiceError(err *errors.Error) *errors.Error
		ServiceErrorWithFields(err *errors.Error, fields ...zap.Field) *errors.Error
		ServiceTxError(err error) *errors.Error
		ServiceDatabaseError(err error) *errors.Error
		ServiceGrpcAdapterError(err error) *errors.Error
		ServiceBrokerAdapterError(err error) *errors.Error

		Info(msg string, fields ...zap.Field)
	}

	logger struct {
		log *zap.Logger
	}
)

func NewLogger(log *zap.Logger) Logger {
	return &logger{
		log: log,
	}
}

func (l *logger) Zap() *zap.Logger {
	return l.log
}

func (l *logger) Sync() error {
	return l.log.Sync()
}

func (l *logger) Panic(text string, fields ...zap.Field) {
	l.log.Panic(text, fields...)
}

func (l *logger) Info(text string, fields ...zap.Field) {
	l.log.Info(text, fields...)
}

func (l *logger) Warn(text string, fields ...zap.Field) {
	l.log.Warn(text, fields...)
}

func (l *logger) Named(name string) Logger {
	return &logger{
		log: l.log.Named(name),
	}
}

func (l *logger) Error(err, reason error) error {
	if reason == nil {
		reason = err
	}

	l.log.Error(formError(err, reason.Error()))
	return err
}

func (l *logger) ErrorRepo(err, reason error, query string) error {
	if reason == nil {
		reason = err
	}

	e, zapMethodField, zapReasonField, zapStackField, zapCallerField, zapFullStackField := formError(err, reason.Error())
	errorMsg := fmt.Sprintf("%s \n %s", e, query)
	l.log.Error(errorMsg, zapMethodField, zapReasonField, zapStackField, zapCallerField, zapFullStackField)
	return err
}

func (l *logger) ServiceError(err *errors.Error) *errors.Error {
	e := fmt.Errorf(err.Reason)
	if err.Details != nil {
		e = err.Details
	}

	l.log.Error(formError(e, err.Reason))
	return err
}

func (l *logger) ServiceGrpcAdapterError(err error) *errors.Error {
	return l.ServiceError(errors.GrpcError(err))
}

func (l *logger) ServiceBrokerAdapterError(err error) *errors.Error {
	return l.ServiceError(errors.BrokerError(err))
}

func (l *logger) ServiceErrorWithFields(err *errors.Error, fields ...zap.Field) *errors.Error {
	e := fmt.Errorf(err.Reason)
	if err.Details != nil {
		e = err.Details
	}

	errData, method, reason, stack, caller, fullStack := formError(e, err.Reason)
	fieldList := []zap.Field{method, reason, stack, caller, fullStack}
	fieldList = append(fieldList, fields...)
	l.log.Error(errData, fieldList...)
	return err
}

func (l *logger) ServiceErrorWithReason(err *errors.Error, reason string) *errors.Error {
	l.log.Error(formError(err.Details, reason))
	return err
}

func (l *logger) ServiceTxError(err error) *errors.Error {
	l.log.Error(formError(err, repository_errors.PgTx.Error()))
	return errors.DatabaseError(err)
}

func (l *logger) ServiceDatabaseError(err error) *errors.Error {
	return l.ServiceError(errors.DatabaseError(err))
}
