package jot

import (
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	text1 := `this is some text,
there are lots of lines
some are the same as text2
but one is not`
	text2 := `this is some text
there are lots of lines
maybe even more!
some are the same as text2
but one is not`
	merged := MergeText(text1, text2)
	mergedTrue := `this is some text,
this is some text
there are lots of lines
maybe even more!
some are the same as text2
but one is not
`

	if mergedTrue != merged {
		t.Errorf("Incorrect merge: %s", merged)
	}

}

func TestEncryptedMerge(t *testing.T) {
	text1 := `this is some text,
there are lots of lines
some are the same as text2
but one is not`
	text2 := `this is some text
there are lots of lines
maybe even more!
some are the same as text2
but one is not`
	mergedTrue := `this is some text,
this is some text
there are lots of lines
maybe even more!
some are the same as text2
but one is not
`
	text := `-----BEGIN PGP SIGNATURE-----

<<<<<<< HEAD
` + StrExtract(EncryptString(text1, Passphrase), "SIGNATURE-----", "-----END", 1) + `
=======
` + StrExtract(EncryptString(text2, Passphrase), "SIGNATURE-----", "-----END", 1) + `
>>>>>>> c85515718f6d26f2279b7a370828a0fc77f16cd8
-----END PGP SIGNATURE-----`
	merged := MergeEncrypted(text, Passphrase)
	if mergedTrue != merged {
		t.Errorf("Incorrect merge: %s", merged)
	}

}

// The following code comes from https://github.com/aryann/difflib/blob/5561ce058dd15606a81ccef3e9b9109dd11036eb/difflib_test.go
// Copyright 2012 Aryan Naraghi (aryan.naraghi@gmail.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var lcsTests = []struct {
	seq1 string
	seq2 string
	lcs  int
}{
	{"", "", 0},
	{"abc", "abc", 3},
	{"mzjawxu", "xmjyauz", 4},
	{"human", "chimpanzee", 4},
	{"Hello, world!", "Hello, world!", 13},
	{"Hello, world!", "H     e    l  l o ,   w  o r l  d   !", 13},
}

func TestLongestCommonSubsequenceMatrix(t *testing.T) {
	for i, test := range lcsTests {
		seq1 := strings.Split(test.seq1, "")
		seq2 := strings.Split(test.seq2, "")
		matrix := longestCommonSubsequenceMatrix(seq1, seq2)
		lcs := matrix[len(matrix)-1][len(matrix[0])-1] // Grabs the lower, right value.
		if lcs != test.lcs {
			t.Errorf("%d. longestCommonSubsequence(%v, %v)[last][last] => %d, expected %d",
				i, seq1, seq2, lcs, test.lcs)
		}
	}
}

var numEqualStartAndEndElementsTests = []struct {
	seq1  string
	seq2  string
	start int
	end   int
}{
	{"", "", 0, 0},
	{"abc", "", 0, 0},
	{"", "abc", 0, 0},
	{"abc", "abc", 3, 0},
	{"abhelloc", "abbyec", 2, 1},
	{"abchello", "abcbye", 3, 0},
	{"helloabc", "byeabc", 0, 3},
}

func TestNumEqualStartAndEndElements(t *testing.T) {
	for i, test := range numEqualStartAndEndElementsTests {
		seq1 := strings.Split(test.seq1, "")
		seq2 := strings.Split(test.seq2, "")
		start, end := numEqualStartAndEndElements(seq1, seq2)
		if start != test.start || end != test.end {
			t.Errorf("%d. numEqualStartAndEndElements(%v, %v) => (%d, %d), expected (%d, %d)",
				i, seq1, seq2, start, end, test.start, test.end)
		}
	}
}
