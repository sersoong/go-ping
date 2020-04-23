package utils

import "testing"

func TestPing(t *testing.T) {
	var ret Ret
	ret = Ping("192.168.6.120", 32, 10000)
	t.Log(ret)
	ret = Ping("192.168.6.1", 32, 10000)
	t.Log(ret)
}
