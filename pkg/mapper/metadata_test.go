package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemMetadataView_ReadOnly_Mapper(t *testing.T) {
	md := Metadata{
		previousVertex: "pv",
		systemMetadata: SystemMetadata{"g": {"k": []byte("v")}},
		userMetadata:   UserMetadata{"ug": {"uk": []byte("uv")}},
	}

	view := md.SystemMetadataView()
	assert.Equal(t, []string{"g"}, view.Groups())
	assert.ElementsMatch(t, []string{"k"}, view.Keys("g"))
	assert.Equal(t, []byte("v"), view.Get("g", "k"))

	// verify user metadata methods mutate as expected
	md.AppendKVString("ug", "k2", "v2")
	assert.Equal(t, []byte("v2"), md.userMetadata["ug"]["k2"])
	md.RemoveKey("ug", "k2")
	_, exists := md.userMetadata["ug"]["k2"]
	assert.False(t, exists)
}
