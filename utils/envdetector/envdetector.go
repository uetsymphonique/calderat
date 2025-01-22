package envdetector

import (
	"calderat/utils/logger"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Environment holds information about the current environment
type Environment struct {
	OS              string   // Operating System
	Arch            string   // Architecture (e.g., amd64, arm64)
	CurrentShell    string   // Current shell in use
	AvailableShells []string // Available shells
	ShortnameShells []string
	NetworkInfo     []NetworkDetail // Network interface details
	Logger          *logger.Logger
}

// NetworkDetail holds details about a single network interface
type NetworkDetail struct {
	Name        string   // Interface name
	IPAddresses []string // IP addresses assigned to the interface
	MTU         int      // Maximum Transmission Unit
	Flags       string   // Interface flags (e.g., up, broadcast, etc.)
}

// DetectEnvironment detects and returns the current environment details
func DetectEnvironment(log *logger.Logger) (*Environment, error) {
	env := &Environment{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Logger: log,
	}

	// Detect the current shell
	currentShell, err := detectCurrentShell()
	if err != nil {
		env.Logger.Log(logger.ERROR, "Failed to detect current shell %v", err)
		return nil, err
	}
	env.CurrentShell = currentShell

	// Detect available shells
	availableShells, err := detectAvailableShells()
	if err != nil {
		env.Logger.Log(logger.ERROR, "Failed to detect available shells %v", err)
		return nil, err
	}
	env.AvailableShells = availableShells
	env.ShortnameShells = extractShortnameShells(env.AvailableShells, env.OS)

	// Detect network interfaces
	networkInfo, err := detectNetworkInterfaces()
	if err != nil {
		env.Logger.Log(logger.ERROR, "Failed to detect network interfaces %v", err)
		return nil, err
	}
	env.NetworkInfo = networkInfo

	return env, nil
}

// detectCurrentShell detects the shell currently in use
func detectCurrentShell() (string, error) {
	// Check SHELL environment variable (common on Unix systems)
	shell := os.Getenv("SHELL")
	if shell != "" {
		return shell, nil
	}

	// Check ComSpec environment variable (Windows command shell)
	if runtime.GOOS == "windows" {
		shell = os.Getenv("ComSpec")
		if shell != "" {
			return shell, nil
		}
	}

	// Fallback: Check the process name
	processName, err := os.Executable()
	if err != nil {
		return "", err
	}
	return processName, nil
}

// detectAvailableShells lists the available shells by checking common shell paths
func detectAvailableShells() ([]string, error) {
	shellPaths := []string{
		"/bin/bash",                      // Bash
		"/bin/zsh",                       // Zsh
		"/usr/bin/fish",                  // Fish shell
		"/bin/sh",                        // Default shell
		"/usr/bin/pwsh",                  // PowerShell Core
		"/usr/bin/nu",                    // Nushell
		"C:\\Windows\\System32\\cmd.exe", // Windows Command Prompt
		"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", // PowerShell
	}

	var available []string
	for _, path := range shellPaths {
		if _, err := exec.LookPath(path); err == nil {
			available = append(available, path)
		}
	}
	return available, nil
}

func extractShortnameShells(shellPaths []string, os string) []string {
	shortnames := []string{}
	for _, path := range shellPaths {
		if strings.Contains(strings.ToLower(path), "powershell") && os == "windows" {
			shortnames = append(shortnames, "psh")
		}
		if strings.Contains(strings.ToLower(path), "cmd") && os == "windows" {
			shortnames = append(shortnames, "cmd")
		}
		if strings.Contains(strings.ToLower(path), "sh") && os == "linux" {
			shortnames = append(shortnames, "sh")
		}
	}
	return shortnames
}
func RemoveDuplicates[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, value := range input {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}

	return result
}

// detectNetworkInterfaces lists all network interfaces and their details
func detectNetworkInterfaces() ([]NetworkDetail, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var details []NetworkDetail
	for _, iface := range interfaces {
		var ipAddresses []string

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for interface %s: %w", iface.Name, err)
		}

		for _, addr := range addrs {
			ipAddresses = append(ipAddresses, addr.String())
		}

		details = append(details, NetworkDetail{
			Name:        iface.Name,
			IPAddresses: ipAddresses,
			MTU:         iface.MTU,
			Flags:       iface.Flags.String(),
		})
	}

	return details, nil
}

// getInterfaceAddresses retrieves IP addresses (IPv4 and IPv6) for a given network interface
func (env *Environment) GetAllIPAddresses() ([]string, error) {
	var ipAddresses []string

	for _, info := range env.NetworkInfo {
		for _, addr := range info.IPAddresses {
			// Parse and append only the IP address (exclude CIDR notation)
			ip, _, err := net.ParseCIDR(addr)
			if err != nil {
				return nil, fmt.Errorf("error parsing address %s for interface %s: %w", addr, info.Name, err)
			}
			if ip.To4() != nil {
				ipAddresses = append(ipAddresses, ip.String())
			}
		}
	}

	return ipAddresses, nil
}
