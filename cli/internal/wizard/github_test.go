package wizard

import "testing"

func TestParseGitRemoteSSH(t *testing.T) {
	owner, repo, err := parseGitRemote("git@github.com:vinceglb/releasekit-ios.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "vinceglb" {
		t.Errorf("expected owner 'vinceglb', got %q", owner)
	}
	if repo != "releasekit-ios" {
		t.Errorf("expected repo 'releasekit-ios', got %q", repo)
	}
}

func TestParseGitRemoteHTTPS(t *testing.T) {
	owner, repo, err := parseGitRemote("https://github.com/vinceglb/releasekit-ios.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "vinceglb" {
		t.Errorf("expected owner 'vinceglb', got %q", owner)
	}
	if repo != "releasekit-ios" {
		t.Errorf("expected repo 'releasekit-ios', got %q", repo)
	}
}

func TestParseGitRemoteHTTPSNoGitSuffix(t *testing.T) {
	owner, repo, err := parseGitRemote("https://github.com/vinceglb/releasekit-ios")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "vinceglb" {
		t.Errorf("expected owner 'vinceglb', got %q", owner)
	}
	if repo != "releasekit-ios" {
		t.Errorf("expected repo 'releasekit-ios', got %q", repo)
	}
}

func TestParseGitRemoteHTTP(t *testing.T) {
	owner, repo, err := parseGitRemote("http://github.com/acme/myapp.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "acme" {
		t.Errorf("expected owner 'acme', got %q", owner)
	}
	if repo != "myapp" {
		t.Errorf("expected repo 'myapp', got %q", repo)
	}
}

func TestParseGitRemoteSSHWithTrailingNewline(t *testing.T) {
	owner, repo, err := parseGitRemote("git@github.com:acme/myapp.git\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "acme" {
		t.Errorf("expected owner 'acme', got %q", owner)
	}
	if repo != "myapp" {
		t.Errorf("expected repo 'myapp', got %q", repo)
	}
}

func TestParseGitRemoteInvalid(t *testing.T) {
	_, _, err := parseGitRemote("not-a-remote-url")
	if err == nil {
		t.Error("expected error for invalid remote URL")
	}
}
