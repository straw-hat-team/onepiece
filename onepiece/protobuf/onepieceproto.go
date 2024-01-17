package protobuf

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FullName protoreflect.FullName

func (s FullName) String() string {
	return string(s)
}

func PackageName[Message protoreflect.ProtoMessage](msg Message) FullName {
	return FullName(msg.ProtoReflect().Descriptor().ParentFile().Package())
}

func MessageFullName[Message protoreflect.ProtoMessage](msg Message) FullName {
	return FullName(msg.ProtoReflect().Descriptor().FullName())
}

func StreamID[Message protoreflect.ProtoMessage](message Message, id string) string {
	return fmt.Sprintf("%s.%s", PackageName(message), id)
}
