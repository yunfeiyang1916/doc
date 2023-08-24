// 回文，就是一个字符串从前向后和从后向前读都是一样的
package word

import "unicode"

// 是否是回文
func IsPalindrome(s string) bool {
	//下面这种只支持ascii编码的
	//length := len(s)
	//for i := 0; i < len(s)/2; i++ {
	//	if s[i] != s[length-1-i] {
	//		return false
	//	}
	//}
	if s == "" {
		return true
	}
	//写个支持unicode并且忽略大小写及空格的
	//先分配切片内存，省的append时多次分配内存，提升性能
	var letters = make([]rune, 0, len(s))
	for _, r := range s {
		//是否是字符而不是空格，如果是，则忽略大小写
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	n := len(letters)
	for i := 0; i < n/2; i++ {
		if letters[i] != letters[n-i-1] {
			return false
		}
	}
	return true
}
