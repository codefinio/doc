package core

import "errors"

var (
	ErrUnsupportedEncodingParameters = errors.New("unsupported encoding parameters")
	ErrNoCCITTFaxDecode              = errors.New("CCITTFaxDecode encoding is not yet implemented")
	ErrNoJBIG2Decode                 = errors.New("JBIG2Decode encoding is not yet implemented")
	ErrNoJPXDecode                   = errors.New("JPXDecode encoding is not yet implemented")
	ErrNoPdfVersion                  = errors.New("version not found")
	ErrTypeError                     = errors.New("type check error")
	ErrRangeError                    = errors.New("range check error")
	ErrNotSupported                  = errors.New("feature not currently supported")
	ErrNotANumber                    = errors.New("not a number")
)
