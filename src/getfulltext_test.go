package sdees

import "testing"

func TestGetTextOfOne(t *testing.T) {
	text, err := GetTextOfOne("./gittest10", EncryptOTP("1"), ShortEncrypt("test.txt"))
	if len(text) == 0 || err != nil {
		t.Errorf("Got no text! Or error: " + err.Error())
	}

	text, err = GetTextOfOne("./gittest10", EncryptOTP("76868761"), ShortEncrypt("test.txt"))
	if err == nil {
		t.Errorf("Wrong branch should throw error")
	}

	text, err = GetTextOfOne("./gittest10", EncryptOTP("1"), ShortEncrypt("kjllkjklkjlkj.txt"))
	if err == nil {
		t.Errorf("Wrong document should throw error")
	}
}
