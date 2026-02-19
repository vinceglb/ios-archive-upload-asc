package wizard

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
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
