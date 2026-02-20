package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
)

const (
	websiteURL      = "https://lark.joshburns.xyz"
	polarOrgID      = "9c530107-ccf7-4462-bc5e-cb245f8143a6"
	licenseCacheTTL = 24 * time.Hour
	configDirName   = "lark"
	licenseFileName = "license"
	cacheFileName   = "license_cache.json"
)

// licenseCache represents the cached license validation result.
type licenseCache struct {
	KeyHash   string `json:"key_hash"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// configDir returns ~/.config/lark, creating it if necessary.
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", configDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create config directory: %w", err)
	}
	return dir, nil
}

// getLicenseKey reads the license key from LARK_LICENSE_KEY env var,
// falling back to ~/.config/lark/license file.
func getLicenseKey() string {
	if key := os.Getenv("LARK_LICENSE_KEY"); key != "" {
		return strings.TrimSpace(key)
	}
	dir, err := configDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(dir, licenseFileName))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// saveLicenseKey writes the key to ~/.config/lark/license.
func saveLicenseKey(key string) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, licenseFileName), []byte(key+"\n"), 0600)
}

// hashKey returns a hex-encoded SHA-256 hash of the key for cache storage.
func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h)
}

// loadCache reads the license cache file.
func loadCache() (*licenseCache, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, cacheFileName))
	if err != nil {
		return nil, err
	}
	var c licenseCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// saveCache writes the license cache file.
func saveCache(c *licenseCache) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, cacheFileName), data, 0600)
}

// masterKeys are license keys that bypass Polar validation entirely.
var masterKeys = map[string]bool{
	"DUSTYBLUE": true,
}

// isMasterKey returns true if the key bypasses Polar validation.
func isMasterKey(key string) bool {
	return masterKeys[key]
}

// getOrgID returns the Polar organization ID from env or const.
func getOrgID() string {
	if id := os.Getenv("POLAR_ORGANIZATION_ID"); id != "" {
		return id
	}
	return polarOrgID
}

// newPolarClient creates an unauthenticated Polar client.
func newPolarClient() *polargo.Polar {
	return polargo.New()
}

// validateLicense calls the Polar license key validation endpoint.
// Returns the status string ("granted", "revoked", etc.) or an error.
func validateLicense(key string) (string, error) {
	client := newPolarClient()

	resp, err := client.CustomerPortal.LicenseKeys.Validate(
		context.Background(),
		components.LicenseKeyValidate{
			Key:            key,
			OrganizationID: getOrgID(),
		},
	)
	if err != nil {
		return "", fmt.Errorf("license validation failed: %w", err)
	}

	return string(resp.ValidatedLicenseKey.Status), nil
}

// activateLicense calls the Polar license key activation endpoint.
func activateLicense(key string) error {
	client := newPolarClient()

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	_, err := client.CustomerPortal.LicenseKeys.Activate(
		context.Background(),
		components.LicenseKeyActivate{
			Key:            key,
			OrganizationID: getOrgID(),
			Label:          hostname,
		},
	)
	if err != nil {
		return fmt.Errorf("license activation failed: %w", err)
	}

	return nil
}

// checkLicense is the main entry point for license verification.
// It reads the key, checks the 24h cache, and validates against Polar if needed.
func checkLicense() error {
	key := getLicenseKey()
	if key == "" {
		return fmt.Errorf("no license key found")
	}

	// Master key bypass
	if isMasterKey(key) {
		return nil
	}

	kh := hashKey(key)

	// Check cache
	if c, err := loadCache(); err == nil {
		if c.KeyHash == kh && time.Since(time.Unix(c.Timestamp, 0)) < licenseCacheTTL {
			if c.Status == "granted" {
				return nil
			}
			return fmt.Errorf("license is %s", c.Status)
		}
	}

	// Validate against Polar API
	status, err := validateLicense(key)
	if err != nil {
		// If API is unreachable, check for a valid cache (even if expired)
		if c, err2 := loadCache(); err2 == nil && c.KeyHash == kh && c.Status == "granted" {
			return nil // allow offline play with previously valid cache
		}
		return fmt.Errorf("unable to verify license: %w", err)
	}

	// Cache the result
	_ = saveCache(&licenseCache{
		KeyHash:   kh,
		Status:    status,
		Timestamp: time.Now().Unix(),
	})

	if status != "granted" {
		return fmt.Errorf("license is %s", status)
	}

	return nil
}

// handleActivateCommand handles the `lark activate <key>` subcommand.
func handleActivateCommand(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "\nUsage: lark activate <license-key>\n\n")
		fmt.Fprintf(os.Stderr, "  Get your license at: %s\n\n", websiteURL)
		os.Exit(1)
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		fmt.Fprintf(os.Stderr, "Error: license key cannot be empty\n")
		os.Exit(1)
	}

	fmt.Print("\n  Saving license key... ")
	if err := saveLicenseKey(key); err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error saving license: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	// Master key — skip Polar API entirely
	if isMasterKey(key) {
		fmt.Println("\n  Master key activated! Run 'lark' to start playing.\n")
		return
	}

	fmt.Print("  Activating license... ")
	if err := activateLicense(key); err != nil {
		fmt.Fprintf(os.Stderr, "\n  Warning: activation failed: %s\n", err)
		fmt.Println("  (This is OK — validation will still work)")
	} else {
		fmt.Println("done")
	}

	fmt.Print("  Validating license... ")
	status, err := validateLicense(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n", err)
		os.Exit(1)
	}

	if status != "granted" {
		fmt.Fprintf(os.Stderr, "\n  Error: license status is '%s'\n", status)
		os.Exit(1)
	}

	// Prime cache
	_ = saveCache(&licenseCache{
		KeyHash:   hashKey(key),
		Status:    status,
		Timestamp: time.Now().Unix(),
	})

	fmt.Println("done")
	fmt.Println("\n  License activated successfully! Run 'lark' to start playing.\n")
}

// handleDeactivateCommand handles the `lark deactivate` subcommand.
func handleDeactivateCommand() {
	dir, err := configDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Remove license file and cache
	os.Remove(filepath.Join(dir, licenseFileName))
	os.Remove(filepath.Join(dir, cacheFileName))

	fmt.Println("\n  License key removed.\n")
}
