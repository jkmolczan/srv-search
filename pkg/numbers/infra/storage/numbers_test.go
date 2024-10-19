package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile
}

func TestNewNumbersStorage_ValidFile(t *testing.T) {
	tmpfile := createTempFile(t, "1\n2\n3")
	defer os.Remove(tmpfile.Name())

	ns, err := NewNumbersStorage(tmpfile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, []int{1, 2, 3}, ns.data)
}

func TestNewNumbersStorage_InvalidFile(t *testing.T) {
	_, err := NewNumbersStorage("nonexistentfile.txt")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestNewNumbersStorage_InvalidNumber(t *testing.T) {
	tmpfile := createTempFile(t, "1\nthree")
	defer os.Remove(tmpfile.Name())

	_, err := NewNumbersStorage(tmpfile.Name())
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestGetSortedCollection(t *testing.T) {
	tmpfile := createTempFile(t, "0\n1\n5\n18\n24\n56\n438593\n768492\n849576\n928372\n989311")
	defer os.Remove(tmpfile.Name())

	ns, err := NewNumbersStorage(tmpfile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{0, 1, 5, 18, 24, 56, 438593, 768492, 849576, 928372, 989311}
	result, err := ns.GetSortedCollection()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, expected, result)
}
