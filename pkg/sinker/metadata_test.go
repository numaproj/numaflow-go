package sinker

import (
	"testing"
)

func TestNewSystemMetadata(t *testing.T) {
	t.Run("with nil map", func(t *testing.T) {
		md := NewSystemMetadata(nil)
		if md == nil {
			t.Fatal("expected non-nil SystemMetadata")
		}
		if md.data == nil {
			t.Error("expected non-nil data map")
		}
		if len(md.data) != 0 {
			t.Errorf("expected empty map, got %d elements", len(md.data))
		}
	})

	t.Run("with existing map", func(t *testing.T) {
		data := map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		}
		md := NewSystemMetadata(data)
		if md == nil {
			t.Fatal("expected non-nil SystemMetadata")
		}
		if len(md.data) != 1 {
			t.Errorf("expected 1 group, got %d", len(md.data))
		}
	})
}

func TestSystemMetadata_Groups(t *testing.T) {
	t.Run("with multiple groups", func(t *testing.T) {
		md := NewSystemMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
			"group2": {"key2": []byte("value2")},
		})
		groups := md.Groups()
		if len(groups) != 2 {
			t.Errorf("expected 2 groups, got %d", len(groups))
		}
	})

	t.Run("with empty metadata", func(t *testing.T) {
		md := NewSystemMetadata(nil)
		groups := md.Groups()
		if len(groups) != 0 {
			t.Errorf("expected 0 groups, got %d", len(groups))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *SystemMetadata
		groups := md.Groups()
		if len(groups) != 0 {
			t.Errorf("expected 0 groups for nil metadata, got %d", len(groups))
		}
	})
}

func TestSystemMetadata_Keys(t *testing.T) {
	t.Run("with multiple keys", func(t *testing.T) {
		md := NewSystemMetadata(map[string]map[string][]byte{
			"group1": {
				"key1": []byte("value1"),
				"key2": []byte("value2"),
			},
		})
		keys := md.Keys("group1")
		if len(keys) != 2 {
			t.Errorf("expected 2 keys, got %d", len(keys))
		}
	})

	t.Run("with non-existent group", func(t *testing.T) {
		md := NewSystemMetadata(nil)
		keys := md.Keys("non-existent")
		if len(keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(keys))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *SystemMetadata
		keys := md.Keys("group1")
		if len(keys) != 0 {
			t.Errorf("expected 0 keys for nil metadata, got %d", len(keys))
		}
	})
}

func TestSystemMetadata_Value(t *testing.T) {
	t.Run("with existing key", func(t *testing.T) {
		md := NewSystemMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		value := md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("with non-existent key", func(t *testing.T) {
		md := NewSystemMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		value := md.Value("group1", "non-existent")
		if value != nil {
			t.Errorf("expected nil value, got %v", value)
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *SystemMetadata
		value := md.Value("group1", "key1")
		if len(value) != 0 {
			t.Errorf("expected empty slice for nil metadata, got %v", value)
		}
	})
}

func TestNewUserMetadata(t *testing.T) {
	t.Run("with nil map", func(t *testing.T) {
		md := NewUserMetadata(nil)
		if md == nil {
			t.Fatal("expected non-nil UserMetadata")
		}
		if md.data == nil {
			t.Error("expected non-nil data map")
		}
		if len(md.data) != 0 {
			t.Errorf("expected empty map, got %d elements", len(md.data))
		}
	})

	t.Run("with existing map", func(t *testing.T) {
		data := map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		}
		md := NewUserMetadata(data)
		if md == nil {
			t.Fatal("expected non-nil UserMetadata")
		}
		if len(md.data) != 1 {
			t.Errorf("expected 1 group, got %d", len(md.data))
		}
	})
}

func TestUserMetadata_Groups(t *testing.T) {
	t.Run("with multiple groups", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
			"group2": {"key2": []byte("value2")},
		})
		groups := md.Groups()
		if len(groups) != 2 {
			t.Errorf("expected 2 groups, got %d", len(groups))
		}
	})

	t.Run("with empty metadata", func(t *testing.T) {
		md := NewUserMetadata(nil)
		groups := md.Groups()
		if len(groups) != 0 {
			t.Errorf("expected 0 groups, got %d", len(groups))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		groups := md.Groups()
		if len(groups) != 0 {
			t.Errorf("expected 0 groups for nil metadata, got %d", len(groups))
		}
	})
}

func TestUserMetadata_Keys(t *testing.T) {
	t.Run("with multiple keys", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {
				"key1": []byte("value1"),
				"key2": []byte("value2"),
			},
		})
		keys := md.Keys("group1")
		if len(keys) != 2 {
			t.Errorf("expected 2 keys, got %d", len(keys))
		}
	})

	t.Run("with non-existent group", func(t *testing.T) {
		md := NewUserMetadata(nil)
		keys := md.Keys("non-existent")
		if len(keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(keys))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		keys := md.Keys("group1")
		if len(keys) != 0 {
			t.Errorf("expected 0 keys for nil metadata, got %d", len(keys))
		}
	})
}

func TestUserMetadata_Value(t *testing.T) {
	t.Run("with existing key", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		value := md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("with non-existent key", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		value := md.Value("group1", "non-existent")
		if value != nil {
			t.Errorf("expected nil value, got %v", value)
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		value := md.Value("group1", "key1")
		if len(value) != 0 {
			t.Errorf("expected empty slice for nil metadata, got %v", value)
		}
	})
}
