package cmd

import (
	"bytes"
	"testing"
)

func TestRootHelp(t *testing.T) {
	command := NewRootCmd()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	command.SetOut(out)
	command.SetErr(errOut)
	command.SetArgs([]string{"--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Contains(out.Bytes(), []byte("releasekit-ios")) {
		t.Fatalf("expected help output to mention releasekit-ios, got: %s", out.String())
	}
}

func TestInvalidCommand(t *testing.T) {
	command := NewRootCmd()
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{"bogus"})

	if err := command.Execute(); err == nil {
		t.Fatal("expected an error for invalid command")
	}
}

func TestWizardHelp(t *testing.T) {
	command := NewRootCmd()
	out := &bytes.Buffer{}
	command.SetOut(out)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{"wizard", "--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Contains(out.Bytes(), []byte("guided setup wizard")) {
		t.Fatalf("expected wizard help output, got: %s", out.String())
	}
}

func TestVersionCommand(t *testing.T) {
	command := NewRootCmd()
	out := &bytes.Buffer{}
	command.SetOut(out)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{"version"})

	if err := command.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Contains(out.Bytes(), []byte("releasekit-ios dev")) {
		t.Fatalf("expected default version output, got: %s", out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("commit: unknown")) {
		t.Fatalf("expected commit output, got: %s", out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("build_date: unknown")) {
		t.Fatalf("expected build date output, got: %s", out.String())
	}
}
