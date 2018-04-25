package clair

import (
	"fmt"
	"io"
)

type ImageError int

const (
	GettingSecretError ImageError = iota + 1
	DecodingConfigDataError
	ParseImageNameError
	GettingManifestError
	BearerTokenRequestError
	BearerTokenResponseError
	UnknownManifestError
	PullingLayersError
	GettingClairAddressError
	ConnectingClairClientError
	PostAncestryError
	GetAncestryError
	VulnerableStatus
	NotVulnerableStatus
)

type imageErrorCode interface {
	Code() ImageError
}

// WithCode annotates err with a new code.
// If err is nil, WithCode returns nil.
func WithCode(err error, code ImageError) error {
	if err == nil {
		return nil
	}
	return &ErrorWithCode{
		cause: err,
		code:  code,
	}
}

type ErrorWithCode struct {
	cause error
	code  ImageError
}

func (w *ErrorWithCode) Error() string    { return w.cause.Error() }
func (w *ErrorWithCode) Cause() error     { return w.cause }
func (w *ErrorWithCode) Code() ImageError { return w.code }

func (w *ErrorWithCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}
