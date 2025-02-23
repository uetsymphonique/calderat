package envdetector_test

import (
	"calderat/utils/envdetector"
	"calderat/utils/logger"
	"net"
	"runtime"
	"testing"
)

// TestDetectEnvironment ensures environment detection runs correctly
func TestDetectEnvironment(t *testing.T) {
	log, err := logger.New("DEBUG")
	if err != nil {
		t.Fatalf("Init log failed: %v", err)
	}
	env, err := envdetector.DetectEnvironment(log)
	if err != nil {
		t.Fatalf("DetectEnvironment failed: %v", err)
	}

	if env.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, env.OS)
	}

	if env.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %s, got %s", runtime.GOARCH, env.Arch)
	}
}

// TestGetAllIPAddresses ensures IP extraction works correctly
func TestGetAllIPAddresses(t *testing.T) {
	log, err := logger.New("DEBUG")
	env, err := envdetector.DetectEnvironment(log)
	if err != nil {
		t.Fatalf("Failed to detect environment: %v", err)
	}

	ips, err := env.GetAllIPAddresses()
	if err != nil {
		t.Fatalf("Failed to get IP addresses: %v", err)
	}

	for _, ip := range ips {
		if net.ParseIP(ip) == nil {
			t.Errorf("Invalid IP detected: %s", ip)
		}
	}
}
