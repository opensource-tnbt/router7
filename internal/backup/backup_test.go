package backup_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"router7/internal/backup"
	"testing"
)

func TestArchive(t *testing.T) {
	tmpin, err := ioutil.TempDir("", "backuptest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpin)

	if err := ioutil.WriteFile(filepath.Join(tmpin, "random.seed"), []byte{0xaa, 0xbb}, 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpin, "dhcp4d"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(tmpin, "dhcp4d", "leases.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := backup.Archive(&buf, tmpin); err != nil {
		t.Fatal(err)
	}

	tmpout, err := ioutil.TempDir("", "backuptest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpout)
	tar := exec.Command("tar", "xzf", "-", "-C", tmpout)
	tar.Stdin = &buf
	tar.Stderr = os.Stderr
	if err := tar.Run(); err != nil {
		t.Fatal(err)
	}

	diff := exec.Command("diff", "-ur", tmpin, tmpout)
	diff.Stdout = os.Stdout
	diff.Stderr = os.Stderr
	if err := diff.Run(); err != nil {
		t.Fatal(err)
	}
}
