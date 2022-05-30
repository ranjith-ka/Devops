/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		name      string
		Scopedurl string
		expected  string
	}{
		{
			name:      "Google",
			Scopedurl: "https://google.com",
			expected:  "200 OK",
		},
		{
			name:      "Joke URL",
			Scopedurl: url,
			expected:  "200 OK",
		},
		{
			name:      "Unknown URL",
			Scopedurl: "https://jhjfhofhet.com",
			expected:  "404 NOK",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := verfiyRandomJoke(tt.Scopedurl)
			if result != tt.expected {
				t.Errorf("expected %s, but got %s", tt.expected, result)
			}
		})
	}
}
