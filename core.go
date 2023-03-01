package pinyin_match

import (
	"fmt"
	"strings"
)

var noTone = parseDict()

//var storage = make(map[string][][]string)

func parseDict() map[string]string {
	parseDict := make(map[string]string)
	for key, val := range simplified {
		tempKey := key
		tempVal := val
		for _, item := range tempVal {
			word := string(item)
			if _, ok := parseDict[word]; !ok {
				parseDict[word] = tempKey
			} else {
				parseDict[word] = parseDict[word] + " " + tempKey
			}
		}
	}
	for key, val := range traditional {
		tempKey := key
		tempVal := val
		for _, item := range tempVal {
			word := string(item)
			if _, ok := parseDict[word]; !ok {
				parseDict[word] = tempKey
			} else {
				parseDict[word] = parseDict[word] + " " + tempKey
			}
		}
	}
	return parseDict
}

// Match 关键字匹配内容
func Match(content, keys string) []int {
	if len(content) == 0 || len(keys) == 0 {
		return nil
	}
	content = strings.ToLower(content)
	keys = strings.ReplaceAll(keys, " ", "")
	keys = strings.ToLower(keys)
	//直接匹配
	keyIndex := strings.Index(content, keys)
	if keyIndex != -1 {
		newContent := content[keyIndex:]
		startIndex := len([]rune(content)) - len([]rune(newContent))
		return []int{startIndex, startIndex + len([]rune(keys)) - 1}
	}
	noPyIndex := getIndex(strings.Split(content, ""), [][]string{strings.Split(keys, "")}, keys)
	if noPyIndex != nil {
		return noPyIndex
	}
	py := getPinyin(content)
	return getIndex(py, GetFullKey(keys), keys)
}

func getIndex(py []string, fullString [][]string, keys string) []int {
	for p := 0; p < len(py); p++ {
		for k := 0; k < len(fullString); k++ {
			key := fullString[k]
			keyLength := len(key)
			extend := keyLength == len([]rune(keys))
			isMatch := true
			i := 0
			preSpaceNum := 0
			spaceNum := 0
			if keyLength <= len(py) {
				for ; i < keyLength; i++ {
					if p+i+preSpaceNum < len(py) && i == 0 && py[p+i+preSpaceNum] == " " {
						preSpaceNum += 1
						i -= 1
					} else {
						if p+i+spaceNum < len(py) && py[p+i+spaceNum] == " " {
							spaceNum += 1
							i -= 1
						} else {
							last := ((p + i + 1) >= len(py)) && ((i + 1) >= keyLength)
							if p+i+spaceNum >= len(py) || !point2point(py[p+i+spaceNum], key[i], last, extend) {
								isMatch = false
								break
							}
						}
					}
				}
				if isMatch {
					return []int{p + preSpaceNum, spaceNum + p + i - 1}
				}
			}
		}
	}
	return nil
}

func point2point(test, key string, last, extend bool) bool {
	if test == "ti" {
		fmt.Println("point2point => ", test, "[", key, "]", last, extend)
	}

	if len(test) == 0 {
		return false
	}
	a := strings.Split(test, " ")
	for _, item := range a {
		if len([]rune(item)) > 0 && extend {
			a = append(a, string([]rune(item)[0]))
		}
	}
	if !last {
		for _, item := range a {
			if item == key {
				return true
			}
		}
		return false
	}
	for _, item := range a {
		if strings.Index(item, key) == 0 {
			return true
		}
	}
	return false
}

// getPinyin 获取拼音
func getPinyin(key string) []string {
	var result []string
	for _, item := range key {
		temp := string(item)
		val, ok := noTone[temp]
		if ok {
			temp = val
		}
		result = append(result, temp)
	}
	return result
}

// GetFullKey 获取输入拼音所有组合（切分 + 首字母）
func GetFullKey(key string) [][]string {
	var result [][]string
	bs := wordBreak(key)
	for _, b := range bs {
		item := strings.Split(b, " ")
		last := len(item) - 1
		if strings.Index(item[last], ",") != -1 {
			keys := strings.Split(item[last], ",")
			for _, key := range keys {
				temp := key
				item = item[0 : len(item)-1]
				item = append(item, temp)
				//str, _ := json.Marshal(item)
				result = append(result, append([]string{}, item...))
			}
		} else {
			result = append(result, item)
		}
	}
	// 首字母简拼匹配
	//if len(result) == 0 || len(result[0]) != len(key) {
	//	result = append(result, strings.Split(key, ""))
	//}
	//storage[key] = result
	return result
}

// wordBreak 输入拼音切分
func wordBreak(key string) []string {
	var result []string
	var solutions []string
	keyLen := len([]rune(key))
	possible := make([]bool, keyLen+1)
	for i := 0; i <= keyLen; i++ {
		possible[i] = true
	}
	getAllSolutions(0, key, &result, &solutions, &possible)
	return solutions
}

// getAllSolutions
func getAllSolutions(start int, s string, result, solutions *[]string, possible *[]bool) {
	sLen := len([]rune(s))
	if start == sLen {
		*solutions = append(*solutions, strings.Join(*result, " "))
		return
	}
	for i := start; i < sLen; i++ {
		piece := string([]rune(s)[start : i+1])
		match := false
		// 最后一个音特殊处理，不用打全
		if isLastPreMatch(piece) && (i+1) >= sLen && (i+1) < len(*possible) {
			if len([]rune(piece)) == 1 {
				*result = append(*result, piece)
			} else {
				var s []string
				for _, item := range allPinyin {
					p := item
					if strings.Index(p, piece) == 0 {
						s = append(s, p)
					}
				}
				*result = append(*result, strings.Join(s, ","))
			}
			match = true
		} else {
			if isAllPinyinInclude(piece) && (i+1) < len(*possible) {
				*result = append(*result, piece)
				match = true
			}
		}
		// 最后一个音不特殊处理
		//if isAllPinyinInclude(piece) && (i+1) < len(*possible) {
		//	*result = append(*result, piece)
		//	match = true
		//}
		if match {
			beforeChange := len(*solutions)
			getAllSolutions(i+1, s, result, solutions, possible)
			if len(*solutions) == beforeChange {
				(*possible)[i+1] = false
			}
			*result = (*result)[0 : len(*result)-1]
		}
	}
}

func isLastPreMatch(str string) bool {
	for _, item := range allPinyin {
		if strings.Index(item, str) == 0 {
			return true
		}
	}
	return false
}

func isAllPinyinInclude(key string) bool {
	for _, item := range allPinyin {
		if item == key {
			return true
		}
	}
	return false
}
