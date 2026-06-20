package net

import (
	"testing"
)

func TestIsIPInSubnet(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		cidr    string
		want    bool
		wantErr bool
	}{
		{
			name:    "IP in standard subnet",
			ip:      "192.168.1.50",
			cidr:    "192.168.1.0/24",
			want:    true,
			wantErr: false,
		},
		{
			name:    "IP outside standard subnet",
			ip:      "192.168.2.50",
			cidr:    "192.168.1.0/24",
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid IP string",
			ip:      "192.168.1.999",
			cidr:    "192.168.1.0/24",
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid CIDR string",
			ip:      "192.168.1.50",
			cidr:    "192.168.1.0/33",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsIPInSubnet(tt.ip, tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsIPInSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsIPInSubnet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
