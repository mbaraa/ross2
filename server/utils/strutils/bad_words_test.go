package strutils

import "testing"

func TestBadWords(t *testing.T) {
	bad := []string{"fuck", "motherfucker", "lol", "ass", "cunt", "noob"}
	count := 0
	for _, bi := range bad {
		if IsBadWord(bi) {
			count++
		}
	}

	if count != 4 {
		t.Fatal("bad words count is wrong")
	}
}
