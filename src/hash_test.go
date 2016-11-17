package gojot

import "testing"

func TestHashing(t *testing.T) {
	if HashString("hi") != HashString("hi") {
		t.Errorf("Hashing not working!")
	}
	if integerHash("hi") != integerHash("hi") {
		t.Errorf("Integer Hashing not working!")
	}
	if GetMD5Hash("hi") != GetMD5Hash("hi") {
		t.Errorf("Integer Hashing not working!")
	}
	if HashWithSalt("hi", "salt") != HashWithSalt("hi", "salt") {
		t.Errorf("Hashing with salt not working!")
	}
	if HashWithSalt("hi", "saltq") == HashWithSalt("hi", "salt") {
		t.Errorf("Hashing with salt not working, 2!")
	}
}

func BenchmarkHash(b *testing.B) {
	testHash := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec dignissim lectus. Vivamus dignissim, elit sed lacinia pretium, sem tellus congue felis, vitae efficitur lectus nunc vitae nisl. Sed luctus lacus vitae orci tristique, eu auctor augue semper. Nullam non eros et augue elementum porta. Sed vel turpis nec urna tempus aliquet. Sed a tempor ligula. Donec vehicula, nunc et suscipit pretium, erat nulla pretium ante, nec luctus velit sem vel massa. Nunc fringilla turpis non euismod finibus. Proin ac euismod mi. Mauris ac lorem sit amet nisi consequat cursus. Praesent aliquam a neque at eleifend. Sed ut sem blandit, faucibus sem non, elementum turpis. Etiam aliquet sit amet felis sit amet malesuada. Nam a dui mattis, faucibus leo eu, porta felis. Sed tristique enim et diam efficitur, at aliquam mauris consectetur. Curabitur rutrum sapien non justo ullamcorper bibendum.`
	for n := 0; n < b.N; n++ {
		HashString(testHash)
	}
}

func BenchmarkHashMD5(b *testing.B) {
	testHash := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec dignissim lectus. Vivamus dignissim, elit sed lacinia pretium, sem tellus congue felis, vitae efficitur lectus nunc vitae nisl. Sed luctus lacus vitae orci tristique, eu auctor augue semper. Nullam non eros et augue elementum porta. Sed vel turpis nec urna tempus aliquet. Sed a tempor ligula. Donec vehicula, nunc et suscipit pretium, erat nulla pretium ante, nec luctus velit sem vel massa. Nunc fringilla turpis non euismod finibus. Proin ac euismod mi. Mauris ac lorem sit amet nisi consequat cursus. Praesent aliquam a neque at eleifend. Sed ut sem blandit, faucibus sem non, elementum turpis. Etiam aliquet sit amet felis sit amet malesuada. Nam a dui mattis, faucibus leo eu, porta felis. Sed tristique enim et diam efficitur, at aliquam mauris consectetur. Curabitur rutrum sapien non justo ullamcorper bibendum.`
	for n := 0; n < b.N; n++ {
		GetMD5Hash(testHash)
	}
}
