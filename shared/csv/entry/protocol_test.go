package entry

import (
	"shared/protobuf/pb"
	"testing"
)

func TestProtocol_GetCmdByProtoName(t *testing.T) {
	CSV.Protocol.GetCmdByProtoName(&pb.S2CYggdrasilTaskInfoChange{})
}
