package ota

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

// UpdateInfo contains metadata about a signed OTA update.
type UpdateInfo struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
	URL      string `json:"url"`
}

// Client handles checking and applying OTA updates for CrewOS.
type Client struct {
	CurrentVersion string
	SentinelURL    string
}

// NewClient creates a new OTA update client.
func NewClient(version, url string) *Client {
	return &Client{
		CurrentVersion: version,
		SentinelURL:    url,
	}
}

// CheckForUpdate queries the Sentinel API for a newer OS/Agent version.
func (c *Client) CheckForUpdate() (*UpdateInfo, error) {
	// Placeholder for actual API call to Sentinel
	return nil, nil 
}

// ApplyUpdate downloads and verifies a signed update package.
func (c *Client) ApplyUpdate(update *UpdateInfo) error {
	fmt.Printf("[OTA] Downloading update %s from %s\n", update.Version, update.URL)
	
	resp, err := http.Get(update.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile := "/tmp/crewos_update.tar.gz"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	hash := sha256.New()
	mw := io.MultiWriter(out, hash)

	if _, err := io.Copy(mw, resp.Body); err != nil {
		return err
	}

	calculatedChecksum := hex.EncodeToString(hash.Sum(nil))
	if calculatedChecksum != update.Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", update.Checksum, calculatedChecksum)
	}

	fmt.Printf("[OTA] Verified update %s. Applying to A/B partition...\n", update.Version)
	// In a real implementation, this would use 'dd' or 'rauc' for A/B switching
	return nil
}
