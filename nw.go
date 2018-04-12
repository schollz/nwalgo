// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nwalgo

var (
	Up   byte = 1
	Left byte = 2
	NW   byte = 3
	None byte = 4
)

func Align(a, b string, match, mismatch, gap int) (alignA, alignB string, score int) {

	aLen := len(a) + 1
	bLen := len(b) + 1

	maxLen := aLen
	if maxLen < bLen {
		maxLen = bLen
	}

	aBytes := make([]byte, 0, maxLen)
	bBytes := make([]byte, 0, maxLen)

	f := make([]int, aLen*bLen)
	pointer := make([]byte, aLen*bLen)

	for i := 1; i < aLen; i++ {
		f[i*bLen] = gap * i
		pointer[i*bLen] = Up
	}
	for j := 1; j < bLen; j++ {
		f[j] = gap * j
		pointer[j] = Left
	}

	pointer[0] = None

	for i := 1; i < aLen; i++ {
		for j := 1; j < bLen; j++ {
			matchMismatch := mismatch
			if a[i-1] == b[j-1] {
				matchMismatch = match
			}
			//(i * bLen) + j
			max := f[((i-1)*bLen)+j-1] + matchMismatch
			hgap := f[((i-1)*bLen)+j] + gap
			vgap := f[(i*bLen)+j-1] + gap

			if hgap > max {
				max = hgap
			}
			if vgap > max {
				max = vgap
			}

			p := NW
			if max == hgap {
				p = Up
			} else if max == vgap {
				p = Left
			}

			pointer[(i*bLen)+j] = p
			f[(i*bLen)+j] = max
		}
	}

	i := aLen - 1
	j := bLen - 1

	score = f[(i*bLen)+j]
	for p := pointer[(i*bLen)+j]; p != None; p = pointer[(i*bLen)+j] {
		if p == NW {
			aBytes = append(aBytes, a[i-1])
			bBytes = append(bBytes, b[j-1])
			i--
			j--
		} else if p == Up {
			aBytes = append(aBytes, a[i-1])
			bBytes = append(bBytes, '-')
			i--
		} else if p == Left {
			aBytes = append(aBytes, '-')
			bBytes = append(bBytes, b[j-1])
			j--
		}
	}

	reverse(aBytes)
	reverse(bBytes)

	return string(aBytes), string(bBytes), score
}

func reverse(a []byte) {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - 1 - i
		a[i], a[j] = a[j], a[i]
	}
}
