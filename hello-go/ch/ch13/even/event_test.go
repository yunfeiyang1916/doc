// 奇偶数包测试程序
package even

import (
	"testing"
)

// 测试偶数
func TestEven(t *testing.T) {
	if !Even(10) {
		t.Log(" 10 must be even!")
		t.Fail()
	}
	if Even(7) {
		t.Log(" 7 is not even!")
		t.Fail()
	}
}

// 测试奇数
func TestOdd(t *testing.T) {
	if !Odd(7) {
		t.Log(" 7 must be odd!")
		t.Fail()
	}
	if Odd(10) {
		t.Log(" 10 is not odd!")
		t.Fail()
	}
}
