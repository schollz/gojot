package sdees

import "testing"

// func TestHashID(t *testing.T) {
// 	fmt.Println(StringToHashID("some kind of string"))
// 	t.Errorf("Not correcting listing files:")
// }

func BenchmarkHashIDToString(b *testing.B) {
	logger.Level(2)
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
