package pinyin_match

import (
	"fmt"
	"sort"
	"testing"
)

func TestMatch(t *testing.T) {
	str := "石石弹(牛)"
	var waitMerge [][]int
	//waitMerge = append(waitMerge, Match(str, "guaguagua"))
	//waitMerge = append(waitMerge, Match(str, "sa"))
	index := Match(str, "s")
	if index != nil {
		waitMerge = append(waitMerge, index)
	}
	//waitMerge = append(waitMerge, Match(str, "shi"))
	waitMerge = rangeMerge(waitMerge)
	fmt.Println(waitMerge)
	if len(waitMerge) == 0 {
		fmt.Println("end")
		return
	}
	contentRune := []rune(str)
	for i := 0; i < len(waitMerge); i++ {
		index := waitMerge[i]
		if i > 0 {
			contentRune = []rune(fmt.Sprintf("%s<em>%s</em>%s", string(contentRune[:index[0]+i*9]), string(contentRune[index[0]+i*9:index[1]+1+i*9]), string(contentRune[index[1]+1+i*9:])))
		} else {
			contentRune = []rune(fmt.Sprintf("%s<em>%s</em>%s", string(contentRune[:index[0]]), string(contentRune[index[0]:index[1]+1]), string(contentRune[index[1]+1:])))
		}
	}
	fmt.Println(string(contentRune))
}

func rangeMerge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i][1] >= intervals[i+1][0] || intervals[i][1]+1 == intervals[i+1][0] {
			intervals[i][1] = max(intervals[i][1], intervals[i+1][1])
			intervals = append(intervals[:i+1], intervals[i+2:]...)
			i--
		}
	}
	return intervals
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
