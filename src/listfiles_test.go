package sdees

import "testing"

func TestListFiles(t *testing.T) {
	files := ListFiles("./gittest10")
	if (files[0] == ("test.txt") && files[1] == ("other.txt")) || (files[1] == ("test.txt") && files[0] == ("other.txt")) {
	} else {
		t.Errorf("Not correcting listing files: %v", files)
	}
}

// func BenchmarkListFiles(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		ListFiles("gittest")
// 	}
// }
