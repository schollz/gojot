package sdees

import "testing"

func TestGetTextOfOne(t *testing.T) {
	text, err := GetTextOfOne("./gittest10", "1", "test.txt")
	if len(text) == 0 || err != nil {
		t.Errorf("Got no text! Or error: " + err.Error())
	}

	text, err = GetTextOfOne("./gittest10", "12930812039", "test.txt")
	if err == nil {
		t.Errorf("Wrong branch should throw error")
	}

	text, err = GetTextOfOne("./gittest10", "1", "asdlfkjasdlkfj")
	if err == nil {
		t.Errorf("Wrong document should throw error")
	}
}
