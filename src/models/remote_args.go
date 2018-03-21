package models

import (
	"strings"
)

//RemoteArgs for storing the properties of the remote url that will be requested whenever a executed
type RemoteArgs struct {
	URL     string    `json:"url" bson:"url"`
	Method  string    `json:"method" bson:"method"`
	Headers MapString `json:"headers" bson:"headers"`
	Body    string    `json:"body" bson:"body"`
}

func (args *RemoteArgs) PreInsert() {
	if args.Headers == nil {
		args.Headers = MapString{}
	}

	args.Method = strings.ToUpper(args.Method)
}

func (args *RemoteArgs) Validate() *Error {
	var errFields ErrorFields

	if len(args.URL) == 0 {
		errFields = append(errFields, NewErrorFieldRequired("args.url"))
	}

	if len(args.Method) == 0 {
		errFields = append(errFields, NewErrorFieldRequired("args.method"))
	}

	if _, ok := HttpMethods[args.Method]; !ok {
		errFields = append(errFields, NewErrorFieldInvalid("args.method"))
	}

	return errFields.GenAppError()
}
