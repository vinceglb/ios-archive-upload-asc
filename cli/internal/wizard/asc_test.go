package wizard

import (
	"testing"
)

func TestParseASCAppsOutputBareArray(t *testing.T) {
	raw := `[{"id":"123456789","attributes":{"name":"My App","bundleId":"com.example.myapp"}}]`
	apps, err := parseASCAppsOutput([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != "123456789" {
		t.Errorf("expected ID 123456789, got %s", apps[0].ID)
	}
	if apps[0].Attributes.Name != "My App" {
		t.Errorf("expected name 'My App', got %s", apps[0].Attributes.Name)
	}
	if apps[0].Attributes.BundleID != "com.example.myapp" {
		t.Errorf("expected bundleId 'com.example.myapp', got %s", apps[0].Attributes.BundleID)
	}
}

func TestParseASCAppsOutputDataEnvelope(t *testing.T) {
	raw := `{"data":[{"id":"987654321","attributes":{"name":"Other App","bundleId":"com.example.other"}}]}`
	apps, err := parseASCAppsOutput([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != "987654321" {
		t.Errorf("expected ID 987654321, got %s", apps[0].ID)
	}
	if apps[0].Attributes.Name != "Other App" {
		t.Errorf("expected name 'Other App', got %s", apps[0].Attributes.Name)
	}
}

func TestParseASCAppsOutputEmptyArray(t *testing.T) {
	raw := `[]`
	apps, err := parseASCAppsOutput([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 0 {
		t.Errorf("expected 0 apps, got %d", len(apps))
	}
}

func TestParseASCAppsOutputEmptyData(t *testing.T) {
	raw := `{"data":[]}`
	apps, err := parseASCAppsOutput([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 0 {
		t.Errorf("expected 0 apps, got %d", len(apps))
	}
}

func TestParseASCAppsOutputMultiple(t *testing.T) {
	raw := `[
		{"id":"1","attributes":{"name":"App One","bundleId":"com.one"}},
		{"id":"2","attributes":{"name":"App Two","bundleId":"com.two"}}
	]`
	apps, err := parseASCAppsOutput([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps, got %d", len(apps))
	}
}
