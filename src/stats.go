package jot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func DisplayStats(fileList []string) {
	data := [][]string{}
	for _, file := range fileList {
		cache, _, _ := UpdateCache(RemoteFolder, EncryptOTP(file), false)
		max := 0
		average := float64(0)
		numEntries := float64(0)
		for _, branch := range cache.Branch {
			numWords := len(GetWordsFromText(branch.Text))
			if numWords > max {
				max = numWords
			}
			average += float64(numWords)
			numEntries += 1
		}
		totalWords := int(average)
		average = average / numEntries
		data = append(data, []string{file, Comma(int64(numEntries)), Comma(int64(totalWords)), Comma(int64(average)), Comma(int64(max))})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Document", "# Entries", "Total # words", "Average # words", "Max # words"})
	for _, v := range data {
		table.Append(v)
	}
	fmt.Printf("\n")
	table.Render() // Send output
	fmt.Printf("\n")
}

// The following is from
// https://github.com/dustin/go-humanize/blob/fef948f2d241bd1fd0631108ecc2c9553bae60bf/comma.go

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
func Comma(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}
