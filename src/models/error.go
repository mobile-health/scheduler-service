package models

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mobile-health/scheduler-service/src/utils"
)

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorFields []*ErrorField

type Error struct {
	ID         string      `json:"id"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Errors     ErrorFields `json:"errors"`
	RequestID  string      `json:"request_id"`
}

func NewError(id string, params map[string]interface{}, statusCode int) *Error {
	return &Error{
		ID:         id,
		StatusCode: statusCode,
		Message:    utils.T(id, params),
	}
}

func NewErrorUnexpected(err error, statusCode int) *Error {
	return NewError("unexpected.app_error", map[string]interface{}{
		"Message": err.Error(),
	}, statusCode)
}

func (err *Error) Render(c *gin.Context) {
	requestID, exist := c.Get("request_id")
	if exist {
		err.RequestID = requestID.(string)
	}
	c.JSON(err.StatusCode, err)
	c.Abort()
}

func (errorFields ErrorFields) GenAppError() *Error {
	if len(errorFields) == 0 {
		return nil
	}

	apperr := NewError("model.error.validation.app_error", nil, 400)
	apperr.Errors = errorFields
	return apperr
}

func (errorFields ErrorFields) String() string {
	d, err := json.Marshal(errorFields)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s", d)
}

func NewErrorFieldRequired(field string) *ErrorField {
	return &ErrorField{
		Field:   field,
		Message: utils.T("model.error.validation.required.app_error"),
	}
}

func NewErrorFieldInvalid(field string) *ErrorField {
	return &ErrorField{
		Field:   field,
		Message: utils.T("model.error.validation.invalid.app_error"),
	}
}
