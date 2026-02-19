package wizard

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func detectWorkspaceCandidate(root string) (string, bool) {
	var matches []string

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			switch name {
			case ".git", "Pods", "Carthage", ".build", "DerivedData", "node_modules", ".swiftpm":
				return fs.SkipDir
			}
			if strings.HasSuffix(name, ".xcworkspace") {
				rel, relErr := filepath.Rel(root, path)
				if relErr == nil {
					matches = append(matches, rel)
				}
				return fs.SkipDir
			}
		}
		return nil
	})

	if len(matches) == 0 {
		return "", false
	}
	sort.Strings(matches)
	return matches[0], true
}

// detectAllWorkspaceCandidates returns all .xcworkspace paths found under root.
func detectAllWorkspaceCandidates(root string) []string {
	var matches []string

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			switch name {
			case ".git", "Pods", "Carthage", ".build", "DerivedData", "node_modules", ".swiftpm":
				return fs.SkipDir
			}
			if strings.HasSuffix(name, ".xcworkspace") {
				rel, relErr := filepath.Rel(root, path)
				if relErr == nil {
					matches = append(matches, rel)
				}
				return fs.SkipDir
			}
		}
		return nil
	})

	sort.Strings(matches)
	return matches
}

// DetectSchemes runs xcodebuild -list on a workspace to discover available schemes.
// Returns nil, nil on any failure (graceful fallback).
func DetectSchemes(workspacePath string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "xcodebuild", "-list", "-workspace", workspacePath, "-json")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil //nolint:nilerr // non-blocking; caller falls back to manual input
	}

	var result struct {
		Workspace struct {
			Schemes []string `json:"schemes"`
		} `json:"workspace"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, nil //nolint:nilerr
	}
	return result.Workspace.Schemes, nil
}

var teamIDRegexp = regexp.MustCompile(`DEVELOPMENT_TEAM\s*=\s*([A-Z0-9]{10})\s*;`)

// DetectTeamID reads the first .xcodeproj/project.pbxproj under root and returns
// the most frequently occurring DEVELOPMENT_TEAM value. Returns "", nil if not found.
func DetectTeamID(workspacePath string) (string, error) {
	// Find the first project.pbxproj relative to the workspace directory.
	workspaceDir := filepath.Dir(workspacePath)

	var pbxprojPath string
	_ = filepath.WalkDir(workspaceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || pbxprojPath != "" {
			return nil
		}
		if !d.IsDir() && d.Name() == "project.pbxproj" && strings.Contains(path, ".xcodeproj") {
			pbxprojPath = path
			return fs.SkipAll
		}
		return nil
	})

	if pbxprojPath == "" {
		return "", nil
	}

	content, err := os.ReadFile(pbxprojPath)
	if err != nil {
		return "", nil //nolint:nilerr
	}

	matches := teamIDRegexp.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return "", nil
	}

	// Return the most frequently occurring team ID.
	counts := make(map[string]int)
	for _, m := range matches {
		counts[string(m[1])]++
	}
	var best string
	var bestCount int
	for id, count := range counts {
		if count > bestCount {
			best = id
			bestCount = count
		}
	}
	return best, nil
}
