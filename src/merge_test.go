package sdees

import "testing"

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
