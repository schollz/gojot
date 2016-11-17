package gojot

import (
	"fmt"
	"testing"
)

func TestSummary(t *testing.T) {
	var cache Cache
	cache.Branch = make(map[string]Entry)
	cache.Branch["1"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: `Anogeissus leiocarpa (African birch; Bambara: ngálǎma) is a tall deciduous tree native to savannas of tropical Africa.[1] It is the sole West African species of the genus Anogeissus, a genus otherwise distributed from tropical central and east Africa through tropical Southeast Asia.[1] A. leiocarpa germinates in the new soils produced by seasonal wetlands and grows at the edges of the rainforest, although not in the rainforest, in the savanna, and along riverbanks forming gallery forests. The tree flowers in the rainy season, from June to October. The seeds, winged samaras, are dispersed by ants.`}
	cache.Branch["2"] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: `The 2010 Nigeria Entertainment Awards was the 5th edition of the ceremony and was held on 18 September 2010. The event took place at BMCC Tribeca Performing Art Center, New York City. Omawunmi and Dagrin led the nomination list with 5 and 4 awards respectively.[1]`}
	cache.Branch["3"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: `Paul William McClellan, born February 3, 1966, in San Mateo, California, was a Major League baseball player for the San Francisco Giants.

McClellan, a graduate of Sequoia High School and the College of San Mateo. He was a first round draft pick by the Giants in 1986. He played his last game with the Giants on October 6, 1991. His MLB career earned run average was 5.26. Later, McClellan joined the now defunct Sonoma County Crushers minor league team which operated between 1995 through 2002.`}
	texts, textBranches, _ := CombineEntries(cache)
	if SummarizeEntries(texts, textBranches) != `1 - Thu, 07 Apr 2005 22:13:13 (95 words):
  Anogeissus leiocarpa (African birch; Bambara: ngálǎma) is a tall deciduous tree
2 - Fri, 08 Apr 2005 22:13:13 (46 words):
  The 2010 Nigeria Entertainment Awards was the 5th edition of the ceremony and
3 - Sat, 09 Apr 2005 22:13:13 (87 words):
  Paul William McClellan, born February 3, 1966, in San Mateo, California, was a` {
		t.Errorf("Incorrect summary")
		fmt.Println(SummarizeEntries(texts, textBranches))
	}
}
