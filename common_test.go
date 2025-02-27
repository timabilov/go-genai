package genai

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetValueByPath(t *testing.T) {
	tests := []struct {
		name  string
		data  map[string]any
		keys  []string
		value any
		want  map[string]any
	}{
		{
			name:  "Simple",
			data:  map[string]any{},
			keys:  []string{"a", "b"},
			value: "v",
			want:  map[string]any{"a": map[string]any{"b": "v"}},
		},
		{
			name:  "Nested",
			data:  map[string]any{"a": map[string]any{}},
			keys:  []string{"a", "b", "c"},
			value: "v",
			want:  map[string]any{"a": map[string]any{"b": map[string]any{"c": "v"}}},
		},
		{
			name:  "String_Array",
			data:  map[string]any{},
			keys:  []string{"b[]", "c"},
			value: []string{"v3", "v4"},
			want:  map[string]any{"b": []map[string]any{{"c": "v3"}, {"c": "v4"}}},
		},
		{
			name:  "Any_Array",
			data:  map[string]any{},
			keys:  []string{"a", "b[]", "c"},
			value: []any{"v1", "v2"},
			want:  map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}}},
		},
		{
			name:  "Array_Existing",
			data:  map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}}},
			keys:  []string{"a", "b[]", "d"},
			value: "v3",
			want:  map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1", "d": "v3"}, {"c": "v2", "d": "v3"}}}},
		},
		{
			name:  "Nil_value",
			data:  map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}}},
			keys:  []string{"a", "b[]", "d"},
			value: nil,
			want:  map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			setValueByPath(tt.data, tt.keys, tt.value)
			if diff := cmp.Diff(tt.data, tt.want); diff != "" {
				t.Errorf("setValueByPath() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestGetValueByPath(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		keys      []string
		want      any
		wantPanic bool
	}{
		{
			name: "Simple",
			data: map[string]any{"a": map[string]any{"b": "v"}},
			keys: []string{"a", "b"},
			want: "v",
		},
		{
			name: "Array_Starting_Element",
			data: map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}},
			keys: []string{"b[]", "c"},
			want: []any{"v1", "v2"},
		},
		{
			name: "Array_Middle_Element",
			data: map[string]any{"a": map[string]any{"b": []map[string]any{{"c": "v1"}, {"c": "v2"}}}},
			keys: []string{"a", "b[]", "c"},
			want: []any{"v1", "v2"},
		},
		{
			name: "KeyNotFound",
			data: map[string]any{"a": map[string]any{"b": "v"}},
			keys: []string{"a", "c"},
			want: nil,
		},
		{
			name: "NilData",
			data: nil,
			keys: []string{"a", "b"},
			want: nil,
		},
		{
			name: "WrongData",
			data: "data",
			keys: []string{"a", "b"},
			want: nil,
		},
		{
			name: "Self",
			data: map[string]any{"a": map[string]any{"b": "v"}},
			keys: []string{"_self"},
			want: map[string]any{"a": map[string]any{"b": "v"}},
		},
		{
			name: "empty key",
			data: map[string]any{"a": map[string]any{"b": "v"}},
			keys: []string{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("The code panicked unexpectedly: %v", r)
					}
				} else {
					if tt.wantPanic {
						t.Errorf("The code did not panic as expected")
					}
				}
			}()

			if tt.wantPanic {
				_ = getValueByPath(tt.data, tt.keys) // This should panic
			} else {
				got := getValueByPath(tt.data, tt.keys)
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("getValueByPath() mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}
