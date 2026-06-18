package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid Configuration",
			cfg: Config{
				Port:        8080,
				SecretKey:   "secure-secret-key-32-chars-long!",
				IPWhitelist: []string{"192.168.1.1", "10.0.0.5"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Port Low",
			cfg: Config{
				Port:        80,
				SecretKey:   "secure-secret-key-32-chars-long!",
				IPWhitelist: []string{"192.168.1.1"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Port High",
			cfg: Config{
				Port:        70000,
				SecretKey:   "secure-secret-key-32-chars-long!",
				IPWhitelist: []string{"192.168.1.1"},
			},
			wantErr: true,
		},
		{
			name: "Secret Key Too Short",
			cfg: Config{
				Port:        8080,
				SecretKey:   "short-key",
				IPWhitelist: []string{"192.168.1.1"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Whitelist IP",
			cfg: Config{
				Port:        8080,
				SecretKey:   "secure-secret-key-32-chars-long!",
				IPWhitelist: []string{"192.168.1.300"}, // Invalid IP octet
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(&tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
