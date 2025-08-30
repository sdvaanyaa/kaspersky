package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "defaults",
			env:  map[string]string{},
			want: Config{Workers: 4, QueueSize: 64},
		},
		{
			name: "valid env",
			env:  map[string]string{"WORKERS": "5", "QUEUE_SIZE": "100"},
			want: Config{Workers: 5, QueueSize: 100},
		},
		{
			name: "invalid env",
			env:  map[string]string{"WORKERS": "xxx", "QUEUE_SIZE": "xxx"},
			want: Config{Workers: 4, QueueSize: 64},
		},
		{
			name: "below min",
			env:  map[string]string{"WORKERS": "0", "QUEUE_SIZE": "-1"},
			want: Config{Workers: 4, QueueSize: 64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := LoadConfig()
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
