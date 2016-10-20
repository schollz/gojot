package sdees

import "testing"

func TestAlliterate(t *testing.T) {
	logger.Debug("Example alliteration: %s", MakeAlliteration())
	if MakeAlliteration() == MakeAlliteration() {
		t.Errorf("Alliterations are not random")
	}
}

func BenchmarkAlliterate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		MakeAlliteration()
	}
}
