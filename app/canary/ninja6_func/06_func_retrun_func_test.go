package main

import (
	"reflect"
	"testing"
)

func Test_incrementor(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"incrementor", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := incrementor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("incrementor() = %v, want %v", got, tt.want)
			}
		})
	}
}
