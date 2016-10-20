package sdees

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/speps/go-hashids"
)

// --------------------------------------------------------
// Some time I did something                 Jan 3rd, 2016
// --------------------------------------------------------
// "Some time I did something" -> B32Encrypted as branch name

func StringToHashID(s string) string {
	sb := []byte(s)
	i := 0
	allInts := []int{}
	for {
		buffer := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		i += 8
		max := 8
		if i > len(sb) {
			max = 8 - (i - len(sb))
		}
		for j := 0; j < max; j++ {
			buffer[j] = sb[i+j-8]
		}
		var num int64
		binary.Read(bytes.NewBuffer(buffer[:]), binary.LittleEndian, &num)
		allInts = append(allInts, int(num))
		// fmt.Printf("\n%v\n%d\n", buffer, num)
		if i > len(sb) {
			break
		}
	}

	hd := hashids.NewData()
	hd.Salt = "a totally cret andom string generated and saved and encrypted using the passphrase"
	hd.MinLength = 30
	h := hashids.NewWithData(hd)
	toEncode := allInts
	e, _ := h.Encode(toEncode)
	logger.Debug("Encoded '%s' as '%s'\n", s, e)

	return string(e)
}

func HashIDToString(e string) string {
	hd := hashids.NewData()
	hd.Salt = "a totally cret andom string generated and saved and encrypted using the passphrase"
	hd.MinLength = 30
	h := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(e)
	var bs []byte
	for _, num := range d {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int64(num))
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		bs = append(bs, buf.Bytes()...)
	}
	// fmt.Printf("\n%v\n'%s'\n\n", bs, string(bs))
	logger.Debug("Decoded '%s' as '%s'\n", e, string(bs))
	return strings.TrimSpace(string(bs))
}
