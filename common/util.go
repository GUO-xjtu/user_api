package common

import (
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"
)

func MarshalJsonPb(message proto.Message) ([]byte, error) {
	if v, ok := message.(*wrappers.BytesValue); ok {
		return v.Value, nil
	}
	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{
		EmitDefaults:          true,
		EnumsAsInts:           true,
		OrigName:              true,
	}
	if err := marshaler.Marshal(&buf, message); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetMSTimeStamp(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
