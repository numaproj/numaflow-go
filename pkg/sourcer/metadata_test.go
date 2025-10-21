package sourcer

import (
	"testing"
)

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

func TestUserMetadata_SetKVGroup(t *testing.T) {
	t.Run("set new group", func(t *testing.T) {
		md := NewUserMetadata(nil)
		md.SetKVGroup("group1", map[string][]byte{"key1": []byte("value1")})
		value := md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("replace existing group", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		md.SetKVGroup("group1", map[string][]byte{"key2": []byte("value2")})
		value := md.Value("group1", "key2")
		if string(value) != "value2" {
			t.Errorf("expected 'value2', got '%s'", string(value))
		}
		// Old key should not exist
		value = md.Value("group1", "key1")
		if value != nil {
			t.Errorf("expected old key to be replaced, got %v", value)
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.SetKVGroup("group1", map[string][]byte{"key1": []byte("value1")})
		// Should not panic
	})
}

func TestUserMetadata_AddKV(t *testing.T) {
	t.Run("add to new group", func(t *testing.T) {
		md := NewUserMetadata(nil)
		md.AddKV("group1", "key1", []byte("value1"))
		value := md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("add to existing group", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		md.AddKV("group1", "key2", []byte("value2"))
		value := md.Value("group1", "key2")
		if string(value) != "value2" {
			t.Errorf("expected 'value2', got '%s'", string(value))
		}
		// Old key should still exist
		value = md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.AddKV("group1", "key1", []byte("value1"))
		// Should not panic
	})
}

func TestUserMetadata_AddKVString(t *testing.T) {
	t.Run("add string value", func(t *testing.T) {
		md := NewUserMetadata(nil)
		md.AddKVString("group1", "key1", "value1")
		value := md.Value("group1", "key1")
		if string(value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", string(value))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.AddKVString("group1", "key1", "value1")
		// Should not panic
	})
}

func TestUserMetadata_AddKVInt(t *testing.T) {
	t.Run("add int value", func(t *testing.T) {
		md := NewUserMetadata(nil)
		md.AddKVInt("group1", "key1", 123)
		value := md.Value("group1", "key1")
		if string(value) != "123" {
			t.Errorf("expected '123', got '%s'", string(value))
		}
	})

	t.Run("with negative int", func(t *testing.T) {
		md := NewUserMetadata(nil)
		md.AddKVInt("group1", "key1", -456)
		value := md.Value("group1", "key1")
		if string(value) != "-456" {
			t.Errorf("expected '-456', got '%s'", string(value))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.AddKVInt("group1", "key1", 123)
		// Should not panic
	})
}

func TestUserMetadata_RemoveKey(t *testing.T) {
	t.Run("remove existing key", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {
				"key1": []byte("value1"),
				"key2": []byte("value2"),
			},
		})
		md.RemoveKey("group1", "key1")
		value := md.Value("group1", "key1")
		if value != nil {
			t.Errorf("expected key to be removed, got %v", value)
		}
		// Other key should still exist
		value = md.Value("group1", "key2")
		if string(value) != "value2" {
			t.Errorf("expected 'value2', got '%s'", string(value))
		}
	})

	t.Run("remove non-existent key", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		md.RemoveKey("group1", "non-existent")
		// Should not panic
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.RemoveKey("group1", "key1")
		// Should not panic
	})
}

func TestUserMetadata_RemoveGroup(t *testing.T) {
	t.Run("remove existing group", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
			"group2": {"key2": []byte("value2")},
		})
		md.RemoveGroup("group1")
		groups := md.Groups()
		if len(groups) != 1 {
			t.Errorf("expected 1 group after removal, got %d", len(groups))
		}
	})

	t.Run("remove non-existent group", func(t *testing.T) {
		md := NewUserMetadata(map[string]map[string][]byte{
			"group1": {"key1": []byte("value1")},
		})
		md.RemoveGroup("non-existent")
		// Should not panic
		groups := md.Groups()
		if len(groups) != 1 {
			t.Errorf("expected 1 group, got %d", len(groups))
		}
	})

	t.Run("with nil metadata", func(t *testing.T) {
		var md *UserMetadata
		md.RemoveGroup("group1")
		// Should not panic
	})
}
