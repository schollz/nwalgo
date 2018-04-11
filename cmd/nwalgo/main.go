// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/schollz/nwalgo"
)

var seq1 = flag.String("seq1", "", "first sequence")
var seq2 = flag.String("seq2", "", "second sequence")
var match = flag.Int("match", 1, "match score")
var mismatch = flag.Int("mismatch", -1, "mismatch score")
var gap = flag.Int("gap", -1, "gap penalty")

type Patch struct {
	Left    string `json:"l"`
	Right   string `json:"r"`
	Between string `json:"b"`
}

func main() {
	flag.Parse()
	if *seq1 == "" || *seq2 == "" {
		log.Fatal("Please provide 2 sequences to align. See nwalgo --help")
	}
	*seq1 = strings.Replace(*seq1, "-", "**dash**", -1)
	*seq2 = strings.Replace(*seq2, "-", "**dash**", -1)
	aln1 := *seq1
	score := 0
	patches := []Patch{}
	for {
		var aln2 string
		aln1, aln2, score = nwalgo.Align(aln1, *seq2, *match, *mismatch, *gap)
		if aln1 == aln2 {
			break
		}
		fmt.Printf("%s\n%s\nScore: %d\n", aln1, aln2, score)
		p := getPatch(aln1, aln2)
		patches = append(patches, p)
		fmt.Println(p)
		aln1 = applyPatch(aln1, p)
		if strings.Replace(aln1, "**dash**", "-", -1) == strings.Replace(aln2, "**dash**", "-", -1) {
			break
		}
	}
	bP, _ := json.MarshalIndent(patches, "", " ")
	fmt.Println(string(bP))
}

func count(s, substr string) int {
	return len(regexp.MustCompile(substr).FindAllStringIndex(s, -1))
}

func applyPatch(s string, p Patch) string {
	pos1 := regexp.MustCompile(p.Left).FindStringIndex(s)[0] + len(p.Left)
	pos2 := regexp.MustCompile(p.Right).FindStringIndex(s)[0]
	return s[:pos1] + p.Between + s[pos2:]
}

func getPatch(aln1, aln2 string) Patch {

	// abcdef
	// ab-def
	//   ^
	bookends := []int{0, 0, 0, 0}
	for i := bookends[0]; i < len(aln1); i++ {
		if aln1[i] != aln2[i] {
			bookends[1] = i
			break
		}
	}
	fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find unique subsequence in front
	for i := bookends[0]; i < bookends[1]; i++ {
		if count(aln1, aln1[i:bookends[1]]) > 1 {
			break
		}
		bookends[0] = i
	}
	fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find where next matching subsequence begins
	bookends[2] = bookends[1]
	for {
		for i := bookends[2]; i < len(aln1); i++ {
			bookends[2] = i
			if aln1[i] == aln2[i] {
				break
			}
		}
		fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[2]:])
		// find where the next matching sequence ends
		bookends[3] = bookends[2]
		for i := bookends[2]; i < len(aln1); i++ {
			if aln1[i] != aln2[i] {
				bookends[3] = i
				break
			}
		}
		if bookends[2] == bookends[3] {
			bookends[3] = len(aln1)
		}
		fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
			break
		}
		bookends[2] = bookends[3]
	}
	// now that we have a second matching sequence, try to reduce it
	for bookends[3] = bookends[2] + 1; bookends[3] < len(aln1); bookends[3]++ {
		fmt.Println(bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
			break
		}
	}

	left := aln1[bookends[0]:bookends[1]]
	left = strings.Replace(left, "**dash**", "-", -1)
	right := aln1[bookends[2]:bookends[3]]
	right = strings.Replace(right, "**dash**", "-", -1)
	insertion := aln2[bookends[1]:bookends[2]]
	insertion = strings.Replace(insertion, "-", "", -1)
	insertion = strings.Replace(insertion, "**dash**", "-", -1)

	fmt.Printf("l: '%s', r: '%s', i: '%s'\n", left, right, insertion)
	return Patch{
		Left:    left,
		Right:   right,
		Between: insertion,
	}
}
