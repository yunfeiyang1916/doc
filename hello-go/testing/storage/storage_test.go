// 用伪邮件发送函数替代真实的邮件发送函数。
package storage

import (
	"strings"
	"testing"
)

func TestCheckQuota(t *testing.T) {
	saved := notifyUser
	//还原更改的通知函数
	defer func() { notifyUser = saved }()
	var u, msg string
	//通知函数更改为记录接收人及信息
	notifyUser = func(user, m string) {
		u, msg = user, m
	}
	const user = "zhangchanglin@360.cn"
	//已使用980MB
	CheckQuota(user)
	if u == "" && msg == "" {
		t.Fatalf("notifyUser not called")
	}
	if u != user {
		t.Errorf("wrong user (%s) notified,want %s", u, user)
	}
	const wantSubstring = "98% of your quota"
	if !strings.Contains(msg, wantSubstring) {
		t.Errorf("unexpected notification message <<%s>>, want substring %q", msg, wantSubstring)
	}
}
