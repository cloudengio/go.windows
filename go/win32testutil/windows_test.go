// Copyright 2021 cloudeng llc. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

//go:build windows

package win32testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInaccessible(t *testing.T) {
	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "test-file.text")
	fatal := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	// File.
	err := os.WriteFile(filename, []byte("hello world\n"), 0600)
	fatal(err)
	_, err = os.ReadFile(filename)
	fatal(err)
	err = MakeInaccessibleToOwner(filename)
	fatal(err)

	_, err = os.ReadFile(filename)
	if err == nil || !strings.Contains(err.Error(), "Access is denied") {
		t.Errorf("missing or incorrect error: %v", err)
	}

	err = os.WriteFile(filename, []byte("hello world\n"), 0600)
	if err == nil || !strings.Contains(err.Error(), "Access is denied") {
		t.Errorf("missing or incorrect error: %v", err)
	}

	err = MakeAccessibleToOwner(filename)
	fatal(err)

	_, err = os.ReadFile(filename)
	fatal(err)

	// Directory.
	dirname := filepath.Join(tmpdir, "test-dir", "sub-dir")
	if err := os.MkdirAll(dirname, 0777); err != nil {
		t.Fatal(err)
	}

	if err := MakeInaccessibleToOwner(dirname); err != nil {
		t.Fatal(err)
	}

	_, err = os.ReadDir(dirname)
	if err == nil || !strings.Contains(err.Error(), "Access is denied") {
		t.Errorf("missing or incorrect error: %v", err)
	}

	err = MakeAccessibleToOwner(dirname)
	fatal(err)
	err = os.ReadDir(dirname)
	fatal(err)
}
