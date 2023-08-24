// 回文测试
package word

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 测试回文
func TestPalindrome(t *testing.T) {
	str1, str2 := "detartrated", "kayak"
	if !IsPalindrome(str1) {
		t.Errorf(`IsPalindrome(%s)=false`, str1)
	}
	if !IsPalindrome(str2) {
		t.Errorf("IsPalindrome(%s)=false", str2)
	}
}

// 测试非回文
func TestNonPalindrome(t *testing.T) {
	str := "palindrome"
	if IsPalindrome(str) {
		t.Errorf("IsPalindrome(%s)=true", str)
	}
}

// 测试中文回文
func TestPalindromeCN(t *testing.T) {
	str := "我是谁谁是我"
	if !IsPalindrome(str) {
		t.Errorf("IsPalindrome(%s)=false", str)
	}
}

// 测试大小写及带空格的回文
func TestCanalPalindrome(t *testing.T) {
	str := "A man, A plan, a canal:Panama"
	if !IsPalindrome(str) {
		t.Errorf("IsPalindrome(%s)=false", str)
	}
}

// 将测试数据合并到表格，这才是最常用的
// 测试失败的信息一般是"f(x)=y,want z",其中f(x)解释了失败的操作和对应的素材，y是实际的运行结果，z是期望的正确的结果。
func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		//输入
		input string
		//期望的正确结果
		want bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, A plan, a canal:Panama", true},
		{"Evil I did dwell: lewd did I live.", true},
		{"我是谁谁是我", true},
		{"呵呵哒啊啊", false},
		{"palindrome", false},
		{"desserts", false},
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q)=%v", test.input, got)
		}
	}
}

// 生成随机回文
func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-i-1] = r
	}
	return string(runes)
}

// 随机测试
func TestRandomPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed:%d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q)=false", p)
		}
	}
}

// 基准测试（性能测试）
func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man, a plan, a canal:Panama")
	}
}

// 示例函数
// 如果示例函数内含有类似下面例子中的注释，go test才会执行这个示例函数，然后检测这个示例函数的标准输出和注释是否匹配
func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("A man, a plan, a canal:Panama"))
	fmt.Println(IsPalindrome("palindrome"))
	//Output:
	//true
	//false
}
