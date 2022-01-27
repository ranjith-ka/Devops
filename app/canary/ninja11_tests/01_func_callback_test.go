package main

import (
	"reflect"
	"testing"
)

func Test_sum(t *testing.T) {
	type args struct {
		xi []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sum test", args{xi: []int{1, 2}}, 3},
		{"sum test", args{xi: []int{1, 4}}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.xi...); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_incrementor(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"incrementor", 1}, // double execute without value to get 1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := incrementor()(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("incrementor() = %v, want %v", got, tt.want)
			}
		})
	}
}
