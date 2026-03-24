package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(time.Minute)
	expected := []byte("pikachu")

	cache.Add("pokemon", expected)

	actual, ok := cache.Get("pokemon")
	if !ok {
		t.Fatal("expected to find key in cache")
	}

	if string(actual) != string(expected) {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestCacheReapsExpiredEntries(t *testing.T) {
	cache := NewCache(20 * time.Millisecond)
	cache.Add("pokemon", []byte("bulbasaur"))

	time.Sleep(50 * time.Millisecond)

	if _, ok := cache.Get("pokemon"); ok {
		t.Fatal("expected entry to be reaped, but it was still found")
	}
}
