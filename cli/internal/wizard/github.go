package wizard

import (
	"fmt"
	"os/exec"
	"strings"
)

// DetectGitRepo detects the GitHub owner and repo from the git remote.
func DetectGitRepo() (owner, repo string, err error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("could not get git remote: %w", err)
	}
	return parseGitRemote(strings.TrimSpace(string(out)))
}

// parseGitRemote parses SSH (git@github.com:owner/repo.git) and
// HTTPS (https://github.com/owner/repo[.git]) GitHub remote URLs.
func parseGitRemote(raw string) (owner, repo string, err error) {
	raw = strings.TrimSpace(raw)

	// SSH: git@github.com:owner/repo.git
	if strings.HasPrefix(raw, "git@") {
		// git@github.com:owner/repo.git
		colonIdx := strings.Index(raw, ":")
		if colonIdx < 0 {
			return "", "", fmt.Errorf("unexpected SSH remote format: %s", raw)
		}
		path := raw[colonIdx+1:]
		path = strings.TrimSuffix(path, ".git")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("unexpected SSH remote path: %s", path)
		}
		return parts[0], parts[1], nil
	}

	// HTTPS: https://github.com/owner/repo[.git]
	if strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "http://") {
		// Strip scheme and host.
		withoutScheme := raw
		if idx := strings.Index(raw, "://"); idx >= 0 {
			withoutScheme = raw[idx+3:]
		}
		// withoutScheme = "github.com/owner/repo[.git]"
		slashIdx := strings.Index(withoutScheme, "/")
		if slashIdx < 0 {
			return "", "", fmt.Errorf("unexpected HTTPS remote format: %s", raw)
		}
		path := withoutScheme[slashIdx+1:]
		path = strings.TrimSuffix(path, ".git")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("unexpected HTTPS remote path: %s", path)
		}
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unrecognized remote URL format: %s", raw)
}

// IsGHAuthenticated returns true if the gh CLI is authenticated.
func IsGHAuthenticated() bool {
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run() == nil
}

// SetGitHubSecrets sets GitHub secrets for the given repo using the gh CLI.
// Returns a map of secret name to error (nil if set successfully).
func SetGitHubSecrets(repoSlug string, secrets map[string]string) map[string]error {
	results := make(map[string]error, len(secrets))
	for name, value := range secrets {
		cmd := exec.Command("gh", "secret", "set",
			"--repo", repoSlug,
			name,
			"--body", value,
		)
		results[name] = cmd.Run()
	}
	return results
}
