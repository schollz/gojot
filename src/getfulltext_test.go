package sdees

import "testing"

func TestGetTextOfOne(t *testing.T) {
	text, err := GetTextOfOne("./gittest10", EncryptOTP("1"), EncryptOTP("test.txt"))
	if len(text) == 0 || err != nil {
		t.Errorf("Got no text! Or error: " + err.Error())
	}

	_, err = GetTextOfOne("./gittest10", EncryptOTP("76868761"), EncryptOTP("test.txt"))
	if err == nil {
		t.Errorf("Wrong branch should throw error")
	}

	_, err = GetTextOfOne("./gittest10", EncryptOTP("1"), EncryptOTP("kjllkjklkjlkj.txt"))
	if err == nil {
		t.Errorf("Wrong document should throw error")
	}
}
