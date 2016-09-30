package sdees

import "testing"

func TestListFiles(t *testing.T) {
	files, _ := ListFiles("./gittest10")
	if (files[0] == "other.txt" || files[1] == "test.txt") && len(files) == 2 {
	} else {
		t.Errorf("Not correcting listing files: %v", files)
	}
}

// func BenchmarkListFiles(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		ListFiles("gittest")
// 	}
// }
