package jot

import (
	"fmt"
	"testing"
)

func TestGenerateAlliteration(t *testing.T) {
	fmt.Printf("\nExample alliteration: %s\n", GenerateAlliteration())
	if GenerateAlliteration() == GenerateAlliteration() {
		t.Errorf("Alliterations are not random")
	}
}

func TestGenerateNonAlliteration(t *testing.T) {
	fmt.Printf("\nExample name: %s\n", GenerateNonAlliteration())
	if GenerateNonAlliteration() == GenerateNonAlliteration() {
		t.Errorf("Alliterations are not random")
	}
}

func BenchmarkAlliterate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateAlliteration()
	}
}
