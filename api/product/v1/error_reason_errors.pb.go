// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsUserNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_USER_NOT_FOUND.String() && e.Code == 404
}

func ErrorUserNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_USER_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsProviderNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_PROVIDER_NOT_FOUND.String() && e.Code == 500
}

func ErrorProviderNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_PROVIDER_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsTokenNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TOKEN_NOT_FOUND.String() && e.Code == 500
}

func ErrorTokenNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_TOKEN_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsSaveProductError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SAVE_PRODUCT_ERROR.String() && e.Code == 500
}

func ErrorSaveProductError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_SAVE_PRODUCT_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsDeleteProductError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DELETE_PRODUCT_ERROR.String() && e.Code == 500
}

func ErrorDeleteProductError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_DELETE_PRODUCT_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsSaveProjectError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SAVE_PROJECT_ERROR.String() && e.Code == 500
}

func ErrorSaveProjectError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_SAVE_PROJECT_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsDeleteProjectError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DELETE_PROJECT_ERROR.String() && e.Code == 500
}

func ErrorDeleteProjectError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_DELETE_PROJECT_ERROR.String(), fmt.Sprintf(format, args...))
}
