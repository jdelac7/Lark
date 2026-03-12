package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// Version is set at build time via ldflags.
var Version = "dev"

const (
	versionURL  = "https://lark.black/releases/version.txt"
	releasesURL = "https://github.com/jdelac7/Lark/releases/latest/download"
)

func handleUpdateCommand() {
	fmt.Printf("\n  Current version: %s\n", Version)
	fmt.Print("  Checking for updates... ")

	// Fetch remote version
	remote, err := fetchRemoteVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	if normalizeVersion(remote) == normalizeVersion(Version) {
		fmt.Println("  Already up to date.")
		fmt.Println()
		return
	}

	fmt.Printf("  New version available: %s\n\n", remote)

	// Determine platform archive
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	archive := fmt.Sprintf("lark_%s_%s.tar.gz", goos, goarch)

	// Download checksum file
	fmt.Print("  Downloading checksums... ")
	checksums, err := fetchText(releasesURL + "/checksums.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	expectedHash := ""
	for _, line := range strings.Split(checksums, "\n") {
		if strings.Contains(line, archive) {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				expectedHash = parts[0]
			}
		}
	}
	if expectedHash == "" {
		fmt.Fprintf(os.Stderr, "  Error: no checksum found for %s\n\n", archive)
		os.Exit(1)
	}

	// Download archive
	fmt.Printf("  Downloading %s... ", archive)
	archiveData, err := fetchBytes(releasesURL + "/" + archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	// Verify checksum
	fmt.Print("  Verifying checksum... ")
	actualHash := fmt.Sprintf("%x", sha256.Sum256(archiveData))
	if actualHash != expectedHash {
		fmt.Fprintf(os.Stderr, "\n  Error: checksum mismatch\n    expected: %s\n    got:      %s\n\n", expectedHash, actualHash)
		os.Exit(1)
	}
	fmt.Println("done")

	// Extract binary from tar.gz
	fmt.Print("  Extracting binary... ")
	binaryData, err := extractBinaryFromTarGz(archiveData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	// Replace current binary
	fmt.Print("  Installing update... ")
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error finding current binary: %s\n\n", err)
		os.Exit(1)
	}

	if err := replaceBinary(execPath, binaryData); err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("done")

	fmt.Printf("\n  Updated to %s\n\n", remote)
}

func handleVersionCommand() {
	fmt.Printf("lark %s\n", Version)
}

func fetchRemoteVersion() (string, error) {
	text, err := fetchText(versionURL)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func fetchText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}
	return string(data), nil
}

func fetchBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

func extractBinaryFromTarGz(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decompressing: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading archive: %w", err)
		}
		if hdr.Name == "lark" && hdr.Typeflag == tar.TypeReg {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("binary 'lark' not found in archive")
}

// replaceBinary atomically replaces the binary at path.
// Writes to a temp file next to the original, then renames.
func replaceBinary(path string, data []byte) error {
	tmpPath := path + ".update"

	if err := os.WriteFile(tmpPath, data, 0755); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("replacing binary: %w", err)
	}

	return nil
}

// normalizeVersion extracts the semver core from a version string.
// "v1.0.0-34-g0c180be" -> "1.0.0", "1.0.0" -> "1.0.0", "v1.0.0" -> "1.0.0"
func normalizeVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	if idx := strings.IndexByte(v, '-'); idx >= 0 {
		v = v[:idx]
	}
	return v
}
