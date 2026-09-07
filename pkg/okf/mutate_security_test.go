package okf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSaveConceptRejectsTraversal verifies that concept paths containing ".."
// segments or absolute paths are refused instead of writing outside the bundle directory.
func TestSaveConceptRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "knowledge")
	if err := InitBundle(bundle); err != nil {
		t.Fatalf("InitBundle: %v", err)
	}

	cases := []struct {
		name string
		path string
	}{
		{"direct parent escape", "../evil.md"},
		{"nested escape", "a/b/../../../evil.md"},
		{"escape via trailing dots", "project/../../evil.md"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Concept{
				ID:    strings.TrimSuffix(tc.path, ".md"),
				Path:  tc.path,
				Type:  "Fact",
				Title: "Evil",
				Body:  "should never be written",
			}
			if err := SaveConcept(bundle, c, true, true, true, "test"); err == nil {
				t.Fatalf("SaveConcept(%q): expected traversal error, got nil", tc.path)
			}
		})
	}

	// Nothing may have been written next to the bundle directory.
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if e.Name() != "knowledge" {
			t.Errorf("unexpected entry created outside bundle: %s", e.Name())
		}
	}
}

// TestSaveConceptAllowsNestedConcept verifies the containment check does not
// reject legitimate nested concept paths, and that bookkeeping stays inside
// the bundle.
func TestSaveConceptAllowsNestedConcept(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "knowledge")
	if err := InitBundle(bundle); err != nil {
		t.Fatalf("InitBundle: %v", err)
	}

	c := &Concept{
		ID:          "project/overview",
		Path:        "project/overview.md",
		Type:        "Fact",
		Title:       "Overview",
		Description: "Project overview.",
		Body:        "Some body text.",
	}
	if err := SaveConcept(bundle, c, true, true, true, "test"); err != nil {
		t.Fatalf("SaveConcept: %v", err)
	}

	if _, err := os.Stat(filepath.Join(bundle, "project", "overview.md")); err != nil {
		t.Errorf("concept file missing inside bundle: %v", err)
	}
	if _, err := os.Stat(filepath.Join(bundle, "project", "index.md")); err != nil {
		t.Errorf("parent index.md missing inside bundle: %v", err)
	}
	if _, err := os.Stat(filepath.Join(bundle, "log.md")); err != nil {
		t.Errorf("log.md missing inside bundle: %v", err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if e.Name() != "knowledge" {
			t.Errorf("unexpected entry created outside bundle: %s", e.Name())
		}
	}
}

// TestSaveConceptRejectsReservedRootFiles ensures index.md and log.md cannot be
// overwritten as concept documents.
func TestSaveConceptRejectsReservedRootFiles(t *testing.T) {
	bundle := t.TempDir()
	if err := InitBundle(bundle); err != nil {
		t.Fatalf("InitBundle: %v", err)
	}

	for _, reserved := range []string{"index.md", "log.md", "."} {
		c := &Concept{
			ID:    strings.TrimSuffix(reserved, ".md"),
			Path:  reserved,
			Type:  "Fact",
			Title: "Reserved Overwrite",
		}
		if err := SaveConcept(bundle, c, true, false, false, "test"); err == nil {
			t.Fatalf("SaveConcept(%q): expected error when targeting reserved root file, got nil", reserved)
		}
	}
}

// TestUpdateParentIndexRejectsTraversal ensures UpdateParentIndex does not write index.md
// outside the bundle directory when given an escaping path.
func TestUpdateParentIndexRejectsTraversal(t *testing.T) {
	bundle := t.TempDir()
	if err := InitBundle(bundle); err != nil {
		t.Fatalf("InitBundle: %v", err)
	}

	c := &Concept{
		ID:   "evil",
		Path: "../evil.md",
	}
	if err := UpdateParentIndex(bundle, c); err == nil {
		t.Fatalf("UpdateParentIndex: expected traversal error, got nil")
	}
}

// TestValidateConceptID verifies upfront concept ID syntax and path containment checks.
func TestValidateConceptID(t *testing.T) {
	invalidIDs := []string{
		"../escaped",
		"../../escaped",
		"/absolute/path",
		"\\windows\\path",
		"sub/../../escaped",
		"index",
		"log",
		"index.md",
		"log.md",
		"",
		".",
		"..",
	}
	for _, id := range invalidIDs {
		if err := ValidateConceptID(id); err == nil {
			t.Errorf("ValidateConceptID(%q) expected error, got nil", id)
		}
	}

	validIDs := []string{
		"architecture/layers",
		"decisions/adr-001",
		"overview",
		"sub/dir/nested-concept",
	}
	for _, id := range validIDs {
		if err := ValidateConceptID(id); err != nil {
			t.Errorf("ValidateConceptID(%q) expected valid, got error: %v", id, err)
		}
	}
}
