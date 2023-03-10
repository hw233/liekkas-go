package strings

//StringDisplayLen 字符串的显示长度 中文=1个字符，英文（包括大小写）=半个字符，数字=半个字符，英文标点符号=半个字符，特殊符号/中文标点符号=1个字符
func StringDisplayLen(s string) float64 {
	sl := 0
	rs := []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			sl++
		} else {
			sl += 2
		}
	}
	return float64(sl) / 2
}
