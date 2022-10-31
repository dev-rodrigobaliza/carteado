package safemap

import (
	"reflect"
	"testing"
)

var (
	safeMap = New[int, int]()
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *SafeMap[int, int]
	}{
		{"int int", safeMap},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New[int, int](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMap_Insert(t *testing.T) {
	tests := []struct {
		name      string
		s         *SafeMap[int, int]
		key       int
		value     int
		wantKey   int
		wantValue int
	}{
		{"int int insert", safeMap, 1, 1, 1, 1},
		{"int int insert with remove", safeMap, 1, 1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Insert(tt.key, tt.value)
			if !tt.s.HasKey(tt.wantKey) {
				t.Errorf("SafeMap.Insert() insert key %v value %v, want key %v", tt.key, tt.value, tt.wantKey)
			}
			value, err := tt.s.GetOneValue(tt.key, false)
			if err != nil {
				t.Errorf("SafeMap.Insert() insert key %v value %v got error %v", tt.key, tt.value, err)
			}
			if value != tt.value {
				t.Errorf("SafeMap.Insert() insert key %v value %v, key %v want value %v", tt.key, tt.value, tt.wantKey, tt.wantValue)
			}
		})
	}
}
