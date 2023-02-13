package pinyin_match

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	str := "石石弹(牛)先曾长还嫋嬝嬲尥褭鳥sadfadsasdfadsfdsfad"
	res := Match(str, "shish")
	fmt.Println(res)
	if res != nil {
		fmt.Println(string([]rune(str)[res[0] : res[1]+1]))
	}
}
