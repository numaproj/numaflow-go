package sourcer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetadataAndUserExposure(t *testing.T) {
	user := UserMetadata{"g1": {"k1": []byte("v1")}}
	md := NewMetadata(user)

	// internal user metadata is the same map (exposed by design for sourcer)
	got := md.UserMetadata()
	assert.Equal(t, []byte("v1"), got["g1"]["k1"])

	// mutate returned map and ensure reflected
	got["g1"]["k2"] = []byte("v2")
	assert.Equal(t, []byte("v2"), md.userMetadata["g1"]["k2"])

	// system metadata view starts empty
	view := md.SystemMetadataView()
	assert.Len(t, view.Groups(), 0)
}

func TestSystemMetadataView_ReadOnly(t *testing.T) {
	// prepare metadata with some system values
	md := Metadata{
		systemMetadata: SystemMetadata{
			"g": {"a": []byte("b")},
		},
		userMetadata: UserMetadata{},
	}

	view := md.SystemMetadataView()
	assert.Equal(t, []string{"g"}, view.Groups())
	assert.ElementsMatch(t, []string{"a"}, view.Keys("g"))
	assert.Equal(t, []byte("b"), view.Get("g", "a"))
}
