package sourcer

import (
	"strconv"
)

type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

// SystemMetadata returns the system metadata.
func (md Metadata) SystemMetadata() SystemMetadata {
	return md.systemMetadata
}

// UserMetadata returns the user metadata.
func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}

// SetKeyValueGroup sets a key-value group to the metadata.
func (md *Metadata) SetKeyValueGroup(group string, kv map[string][]byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group] = kv
}

// AppendKeyValue appends a key-value pair with byte value to the metadata.
func (md *Metadata) AppendKeyValue(group, key string, value []byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = value
}

// AppendKeyValueString appends a key-value pair with string value to the metadata.
func (md *Metadata) AppendKeyValueString(group, key, value string) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = []byte(value)
}

func (md *Metadata) AppendKeyValueInt(group, key string, value int) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = []byte(strconv.Itoa(value))
}
