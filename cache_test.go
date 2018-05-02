package cache

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var (
	testData   = []byte(`{"hello":"world"}`)
	defaultExp = time.Duration(1 * time.Second)
)

func TestFetch(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))

	err := m.Set("greeting", testData, defaultExp)
	if err != nil {
		t.Fatal(err)
	}

	item, err := m.Fetch("greeting")
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(item, testData) > 0 {
		t.Fatalf("expected value: 'world' got: '%s'\n", item)
	}
}

func TestFetchExpired(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))

	err := m.Set("greeting", testData, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = m.Fetch("greeting")
	if err != ErrNotFound {
		t.Fatalf("item expected to be expired")
	}
}

func TestSet(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))
	err := m.Set("greeting", testData, defaultExp)
	if err != nil {
		t.Fatal(err)
	}

	assertStoredItemCountIs(t, m, 1)
}

func TestSetExisting(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))
	err := m.Set("greeting", testData, defaultExp)
	if err != nil {
		t.Fatal(err)
	}

	testData = []byte(`{"hello":"universe"}`)
	err = m.Set("greeting", testData, defaultExp)
	if err != nil {
		t.Fatal(err)
	}

	assertStoredItemCountIs(t, m, 1)
}

func TestSetOverItemLimit(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))
	for i := 0; i < 20; i++ {
		key := fmt.Sprint(i)
		value := []byte(key)
		err := m.Set(key, value, defaultExp)
		if err != nil {
			t.Fatal(err)
		}
	}

	// check count
	assertStoredItemCountIs(t, m, 10)

	// check content
	for i := 10; i < 20; i++ {
		key := fmt.Sprint(i)
		expected := []byte(key)
		v, err := m.Fetch(key)
		if err != nil || bytes.Compare(v, expected) > 0 {
			t.Fatal(err)
		}
	}
}

func TestSetOverBytesLimit(t *testing.T) {
	m := NewMemCache(MaxBytesLimit(8))

	// total of 9 bytes
	testItems := [][]byte{
		[]byte("123"),
		[]byte("456"),
		[]byte("789"),
	}

	for i, item := range testItems {
		key := fmt.Sprint(i)
		err := m.Set(key, item, defaultExp)
		if err != nil {
			t.Fatal(err)
		}
	}

	assertStoredItemCountIs(t, m, 2)
}

func TestDelete(t *testing.T) {
	m := NewMemCache(MaxItemLimit(10))

	err := m.Set("greeting", testData, defaultExp)
	if err != nil {
		t.Fatal(err)
	}

	m.Delete("greeting")
	assertStoredItemCountIs(t, m, 0)
	if m.currentBytes != 0 {
		t.Fatal("currentBytes expected to be 0")
	}
}

func assertStoredItemCountIs(t *testing.T, m *MemCache, expected uint) {
	msg := "expected %v items to be stored, got %v\n"

	if l := uint(m.list.Len()); l != expected {
		t.Fatalf(msg, expected, l)
	}

	if l := uint(len(m.items)); l != expected {
		t.Fatalf(msg, expected, l)
	}
}
