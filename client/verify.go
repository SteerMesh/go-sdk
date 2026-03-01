package client

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

// BundleSignature is the optional signature in a bundle manifest.
type BundleSignature struct {
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"keyId"`
	Value     string `json:"value"`
}

// VerifyBundleManifest reads bundle-manifest.json at manifestPath and, if it contains a signature,
// verifies it with the Ed25519 public key at publicKeyPath (PEM PKIX). No-op if manifest has no signature.
func VerifyBundleManifest(manifestPath, publicKeyPath string) error {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var raw struct {
		Version       string          `json:"version"`
		BundleVersion string          `json:"bundleVersion,omitempty"`
		Packs         []PackRef       `json:"packs"`
		Files         []FileEntry      `json:"files"`
		Signature     *BundleSignature `json:"signature,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw.Signature == nil {
		return nil
	}
	if raw.Signature.Algorithm != "Ed25519" {
		return fmt.Errorf("unsupported signature algorithm: %s", raw.Signature.Algorithm)
	}
	sigBytes, err := base64.StdEncoding.DecodeString(raw.Signature.Value)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	canonical := canonicalManifestBytes(raw.Version, raw.BundleVersion, raw.Packs, raw.Files)
	pemData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("no PEM block in public key file")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	edPub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return fmt.Errorf("key is not Ed25519")
	}
	if !ed25519.Verify(edPub, canonical, sigBytes) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func canonicalManifestBytes(version, bundleVersion string, packs []PackRef, files []FileEntry) []byte {
	c := struct {
		Version       string     `json:"version"`
		BundleVersion string     `json:"bundleVersion,omitempty"`
		Packs         []PackRef  `json:"packs"`
		Files         []FileEntry `json:"files"`
	}{Version: version, BundleVersion: bundleVersion, Packs: packs, Files: files}
	b, _ := json.Marshal(c)
	return b
}
