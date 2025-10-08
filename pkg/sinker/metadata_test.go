package sinker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataViews_Sinker(t *testing.T) {
	md := Metadata{
		previousVertex: "prev",
		systemMetadata: SystemMetadata{"sg": {"sk": []byte("sv")}},
		userMetadata:   UserMetadata{"ug": {"uk": []byte("uv")}},
	}

	// previous vertex
	assert.Equal(t, "prev", md.PreviousVertex())

	// system view
	sv := md.SystemMetadataView()
	assert.Equal(t, []string{"sg"}, sv.Groups())
	assert.ElementsMatch(t, []string{"sk"}, sv.Keys("sg"))
	assert.Equal(t, []byte("sv"), sv.Get("sg", "sk"))

	// user view (read-only)
	uv := md.UserMetadataView()
	assert.Equal(t, []string{"ug"}, uv.Groups())
	assert.ElementsMatch(t, []string{"uk"}, uv.Keys("ug"))
	assert.Equal(t, []byte("uv"), uv.Get("ug", "uk"))
}
