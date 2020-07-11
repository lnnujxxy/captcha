//
// @Author: zhouweiwei
// @Date: 2020/7/11 11:50 上午
//

package captcha

import (
	"bytes"
	"testing"
)

func TestRedisSetGet(t *testing.T) {
	r := NewRedis("localhost:6379", "", WithMaxActive(20), WithMaxIdle(20))
	id := "captcha id"
	d := RandomDigits(10)
	r.Set(id, d)
	d2 := r.Get(id, false)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestRedisGetClear(t *testing.T) {
	r := NewRedis("localhost:6379", "")
	id := "captcha id"
	d := RandomDigits(10)
	r.Set(id, d)
	d2 := r.Get(id, true)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = r.Get(id, false)
	if d2 != nil {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}
