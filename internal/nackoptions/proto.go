package nackoptions

import (
	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
)

// ToProto converts NackOptions to the outgoing proto type.
// It is nil-safe: a nil input returns nil, leaving the wire field unset.
func ToProto(n *NackOptions) *common.NackOptions {
	if n == nil {
		return nil
	}
	return &common.NackOptions{
		Reason:        n.Reason,
		MaxDeliveries: n.MaxDeliveries,
		Delay:         n.Delay,
	}
}

// FromProto converts the incoming proto type to NackOptions.
// It is nil-safe: a nil input returns nil.
func FromProto(p *common.NackOptions) *NackOptions {
	if p == nil {
		return nil
	}
	return &NackOptions{
		Delay:         p.Delay,
		MaxDeliveries: p.MaxDeliveries,
		Reason:        p.Reason,
	}
}
