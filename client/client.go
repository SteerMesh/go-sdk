package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const APIKeyHeader = "X-API-Key"

// Client is the SteerMesh Cloud API client.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// New returns a new API client.
func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

// Pack is pack metadata from the API.
type Pack struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

// ListPacks returns all packs from the registry.
func (c *Client) ListPacks() ([]Pack, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/packs", nil)
	req.Header.Set(APIKeyHeader, c.APIKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("packs: %s %s", resp.Status, string(body))
	}
	var out struct {
		Packs []Pack `json:"packs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Packs, nil
}

// Bundle is a bundle response.
type Bundle struct {
	ID       string       `json:"id"`
	Manifest BundleManifest `json:"manifest"`
	Files    []FileRef     `json:"files"`
}

type BundleManifest struct {
	Version string     `json:"version"`
	Packs   []PackRef  `json:"packs"`
	Files   []FileEntry `json:"files"`
}

type PackRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type FileEntry struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type FileRef struct {
	Path string `json:"path"`
	URL  string `json:"url,omitempty"`
}

// GetBundle returns a bundle by ID.
func (c *Client) GetBundle(bundleID string) (*Bundle, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/bundles/"+bundleID, nil)
	req.Header.Set(APIKeyHeader, c.APIKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bundles/%s: %s %s", bundleID, resp.Status, string(body))
	}
	var b Bundle
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// DownloadBundle fetches the bundle and writes manifest + files to outDir.
func (c *Client) DownloadBundle(bundleID, outDir string) error {
	b, err := c.GetBundle(bundleID)
	if err != nil {
		return err
	}
	// Write manifest
	manifestBytes, _ := json.MarshalIndent(b.Manifest, "", "  ")
	if err := writeFile(outDir, "bundle-manifest.json", manifestBytes); err != nil {
		return err
	}
	for _, f := range b.Files {
		if f.URL == "" {
			continue
		}
		req, _ := http.NewRequest("GET", f.URL, nil)
		req.Header.Set(APIKeyHeader, c.APIKey)
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("download %s: %s", f.Path, resp.Status)
		}
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err := writeFile(outDir, f.Path, data); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(outDir, relPath string, data []byte) error {
	full := filepath.Join(outDir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, data, 0644)
}
