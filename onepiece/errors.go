package onepiece

import "errors"

var ErrTerminalState = errors.New("state is terminal")
var ErrUnknownCommand = errors.New("unknown command")
var ErrUnknownEvent = errors.New("unknown event")
var ErrUnknownMessage = errors.New("unknown message")
