package model

import (
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		wantErrMsg string
	}{
		{
			name:       "empty name is required",
			value:      "",
			wantErrMsg: "name is required",
		},
		{
			name:       "whitespace-only name is required",
			value:      "   ",
			wantErrMsg: "name is required",
		},
		{
			name:  "valid name is accepted",
			value: "web-server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateName(tt.value)

			requireErrorMessage(t, gotErr, tt.wantErrMsg)
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		wantErrMsg string
	}{
		{
			name:       "empty hostname is required",
			value:      "",
			wantErrMsg: "hostname is required",
		},
		{
			name:       "whitespace-only hostname is required",
			value:      "\t\n",
			wantErrMsg: "hostname is required",
		},
		{
			name:  "valid hostname is accepted",
			value: "web-01.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateHostname(tt.value)

			requireErrorMessage(t, gotErr, tt.wantErrMsg)
		})
	}
}

func TestValidateOS(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		wantErrMsg string
	}{
		{
			name:       "empty os is required",
			value:      "",
			wantErrMsg: "os is required",
		},
		{
			name:  "valid os is accepted",
			value: "ubuntu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateOS(tt.value)

			requireErrorMessage(t, gotErr, tt.wantErrMsg)
		})
	}
}

func TestServerValidate(t *testing.T) {
	tests := []struct {
		name       string
		server     *Server
		wantErrMsg string
	}{
		{
			name: "valid server passes validation",
			server: &Server{
				Name:     "web-server",
				Hostname: "web-01.example.com",
				OS:       "ubuntu",
			},
		},
		{
			name: "server with empty name fails validation",
			server: &Server{
				Name:     "",
				Hostname: "web-01.example.com",
				OS:       "ubuntu",
			},
			wantErrMsg: "name is required",
		},
		{
			name: "server with empty hostname fails validation",
			server: &Server{
				Name:     "web-server",
				Hostname: "",
				OS:       "ubuntu",
			},
			wantErrMsg: "hostname is required",
		},
		{
			name: "server with empty os fails validation",
			server: &Server{
				Name:     "web-server",
				Hostname: "web-01.example.com",
				OS:       "",
			},
			wantErrMsg: "os is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.server.Validate()

			requireErrorMessage(t, gotErr, tt.wantErrMsg)
		})
	}
}

// requireErrorMessage fails the test when an unexpected error is returned
// (wantErrMsg is empty) or when the returned error message does not match.
func requireErrorMessage(t *testing.T, gotErr error, wantErrMsg string) {
	t.Helper()

	if wantErrMsg == "" && gotErr != nil {
		t.Fatalf("unexpected error: %v", gotErr)
		return
	}

	if wantErrMsg != "" {
		if gotErr == nil {
			t.Fatalf("expected error %q, got nil", wantErrMsg)
			return
		}

		if !strings.Contains(gotErr.Error(), wantErrMsg) {
			t.Fatalf("error = %q, want it to contain %q", gotErr.Error(), wantErrMsg)
		}
	}
}
