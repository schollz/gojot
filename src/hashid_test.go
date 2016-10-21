package sdees

import "testing"

func TestHashID(t *testing.T) {
	if "some kind of string" != HashIDToString(StringToHashID("some kind of string")) {
		t.Errorf("HashID not working")
	}
}

func BenchmarkHashIDToString(b *testing.B) {
	s := StringToHashID("some kind of string ")
	for n := 0; n < b.N; n++ {
		HashIDToString(s)
	}
}

func BenchmarkStringToHashID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StringToHashID("some kind of string")
	}
}
