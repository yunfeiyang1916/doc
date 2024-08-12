package cert

import (
	"fmt"
	"testing"
)

func TestNewCA(t *testing.T) {
	ca, priv, err := NewCA("test", "test", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(ca.Raw))
	fmt.Println(priv)

}

func TestLoadOrCreateCA(t *testing.T) {
	_, _, err := LoadOrCreateCA("test.key", "test.crt")
	if err != nil {
		t.Fatal(err)
	}
}
