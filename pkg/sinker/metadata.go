package sinker

import "strconv"

// SystemMetadata and UserMetadata are mappings of group name to key-value pairs
type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
	previousVertex string
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

func (md Metadata) PreviousVertex() string {
	return md.previousVertex
}

func (md Metadata) SystemMetadata() SystemMetadata {
	return md.systemMetadata
}

func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}

func (md *Metadata) SetKeyValueGroup(group string, kv map[string][]byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group] = kv
}

func (md *Metadata) AppendKeyValue(group, key string, value []byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = value
}

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
