package onepiecemessage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrMessageTypeInvalid = errors.New("invalid message type")
)

var (
	allowedChars = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[._][a-zA-Z0-9]+)*$`)
)

type MessageType string

func (m MessageType) AsPtr() *MessageType {
	return &m
}

func (m MessageType) String() string {
	return string(m)
}

func NewMessageType(msgType string) (*MessageType, error) {
	if !allowedChars.MatchString(msgType) {
		return nil, fmt.Errorf("%w: package name has invalid characters: %s", ErrMessageTypeInvalid, msgType)
	}

	tokens := strings.Split(msgType, ".")
	if len(tokens) != 5 {
		return nil, fmt.Errorf("%w: package name must be exactly <namespace>.<domain>.<stream>.<stream version>.<message name> tokens: %s", ErrMessageTypeInvalid, msgType)
	}

	return MessageType(msgType).AsPtr(), nil
}
