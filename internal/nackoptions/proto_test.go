package nackoptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func TestToProto_Nil(t *testing.T) {
	assert.Nil(t, ToProto(nil))
}

func TestFromProto_Nil(t *testing.T) {
	assert.Nil(t, FromProto(nil))
}

func TestRoundTrip_AllFields(t *testing.T) {
	in := &NackOptions{
		Delay:         ptr(uint64(500)),
		MaxDeliveries: ptr(uint32(3)),
		Reason:        ptr("retry"),
	}
	p := ToProto(in)
	assert.Equal(t, uint64(500), p.GetDelay())
	assert.Equal(t, uint32(3), p.GetMaxDeliveries())
	assert.Equal(t, "retry", p.GetReason())

	back := FromProto(p)
	assert.Equal(t, in, back)
}

func TestToProto_PartialFields(t *testing.T) {
	in := &NackOptions{Delay: ptr(uint64(100))}
	p := ToProto(in)
	assert.Equal(t, uint64(100), p.GetDelay())
	assert.Nil(t, p.MaxDeliveries)
	assert.Nil(t, p.Reason)
}
