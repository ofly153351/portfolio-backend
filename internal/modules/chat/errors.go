package chat

import "errors"

var ErrNotImplemented = errors.New("chat: not implemented")
var ErrMessageRequired = errors.New("chat: message is required")
var ErrAIServiceURLMissing = errors.New("chat: AI service URL is missing")
