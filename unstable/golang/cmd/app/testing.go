package main

import (
	onepiecepb "buf.build/gen/go/straw-hat-llc/onepiece/protocolbuffers/go/straw-hat-llc/onepiece"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"unstable/plandomain/planproto"
)

func CheckAndCollectStreamIDFields(msg proto.Message) []string {
	s := reflect.ValueOf(msg).Elem()
	msgType := s.Type()
	ids := make([]string, 0)
	for i := 0; i < s.NumField(); i++ {
		fieldDesc := msgType.Field(i)
		if proto.HasExtension(msg, onepiecepb.E_StreamId) {
			ext, err := proto.GetExtension(msg, onepiecepb.E_StreamId)
			if err == nil && *(ext.(*bool)) {
				ids = append(ids, fieldDesc.Name)
			}
		}
	}

	return ids
}

func ProtoTesting() bool {
	m := &planproto.CreatePlan{
		PlanId: "example-plan-id",
	}

	id := CheckAndCollectStreamIDFields(m)
	fmt.Println("Fields with stream_id option:", id)

	return true
}

//
//func getFileDescriptorOfficial(path string) protoreflect.FileDescriptor {
//	buf, err := os.ReadFile("./proto/hello.proto.fds")
//	if err != nil {
//		fmt.Printf("Failed to read FileDescriptorSet file. Error: %s", err)
//		os.Exit(-1)
//	}
//
//	// unmarshal
//	var fds dpb.FileDescriptorSet
//	if err = proto_v2.Unmarshal(buf, &fds); err != nil {
//		fmt.Printf("Failed to load FileDescriptorSet file. Error: %s", err)
//		os.Exit(-1)
//	}
//
//	files, err := protodesc.NewFiles(&fds)
//	if err != nil {
//		fmt.Printf("Failed to new protodesc.Files. Error: %s", err)
//		// ERROR: proto: could not resolve import "google/protobuf/descriptor.proto": not found
//		os.Exit(-1)
//	}
//
//	fd_v2, err := files.FindFileByPath("hello.proto")
//	if err != nil {
//		fmt.Printf("Failed to fild FileDescriptor. Error: %s", err)
//		os.Exit(-2)
//	}
//	return fd_v2
//}
