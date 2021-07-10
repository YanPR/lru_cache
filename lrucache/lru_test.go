package lrucache

import "testing"

func Test_LRUCache(t *testing.T) {
	type setArgs struct {
		key      string
		val      interface{}
		wantFlag int
	}
	type getArgs struct {
		key      string
		wantRes  interface{}
		wantFlag int
	}
	setParas := []setArgs{
		{"1", 1, 2},
		{"2", 2, 2},
		{"1", 1, 1},
		{"3", 3, 2},
	}
	getparas := []getArgs{
		{"1", 1, 1},
		{"2", nil, 0},
		{"4", nil, 0},
	}
	C := New(2)
	for _, p := range setParas {
		if flag := C.Set(p.key, p.val); flag != p.wantFlag {
			t.Errorf("set() = %v, want %v", flag, p.wantFlag)
		}
	}
	for _, p := range getparas {
		if res, flag := C.Get(p.key); res != p.wantRes || flag != p.wantFlag {
			t.Errorf("get() = %v %v, want %v %v", res, flag, p.wantRes, p.wantFlag)
		}
	}
}
