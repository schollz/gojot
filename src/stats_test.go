package jot

import "testing"

// The following is from
// https://github.com/dustin/go-humanize

// Copyright (c) 2005-2008  Dustin Sallings <dustin@spy.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// <http://www.opensource.org/licenses/mit-license.php>
// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Comma(834142) -> 834,142
func TestCommas(t *testing.T) {
	testList{
		{"0", Comma(0), "0"},
		{"10", Comma(10), "10"},
		{"100", Comma(100), "100"},
		{"1,000", Comma(1000), "1,000"},
		{"10,000", Comma(10000), "10,000"},
		{"100,000", Comma(100000), "100,000"},
		{"10,000,000", Comma(10000000), "10,000,000"},
		{"10,100,000", Comma(10100000), "10,100,000"},
		{"10,010,000", Comma(10010000), "10,010,000"},
		{"10,001,000", Comma(10001000), "10,001,000"},
		{"123,456,789", Comma(123456789), "123,456,789"},
		{"maxint", Comma(9.223372e+18), "9,223,372,000,000,000,000"},
		{"minint", Comma(-9.223372e+18), "-9,223,372,000,000,000,000"},
		{"-123,456,789", Comma(-123456789), "-123,456,789"},
		{"-10,100,000", Comma(-10100000), "-10,100,000"},
		{"-10,010,000", Comma(-10010000), "-10,010,000"},
		{"-10,001,000", Comma(-10001000), "-10,001,000"},
		{"-10,000,000", Comma(-10000000), "-10,000,000"},
		{"-100,000", Comma(-100000), "-100,000"},
		{"-10,000", Comma(-10000), "-10,000"},
		{"-1,000", Comma(-1000), "-1,000"},
		{"-100", Comma(-100), "-100"},
		{"-10", Comma(-10), "-10"},
	}.validate(t)
}

type testList []struct {
	name, got, exp string
}

func (tl testList) validate(t *testing.T) {
	for _, test := range tl {
		if test.got != test.exp {
			t.Errorf("On %v, expected '%v', but got '%v'",
				test.name, test.exp, test.got)
		}
	}
}
