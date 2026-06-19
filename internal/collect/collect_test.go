package collect

import (
	"sync/atomic"
	"testing"
)

func TestMemoizeSameKey(t *testing.T) {
	var calls int32
	c := New(1)
	c.run = func(name string, args ...string) string {
		atomic.AddInt32(&calls, 1)
		return "out"
	}
	_ = c.GitRemote()
	_ = c.GitRemote()
	if calls != 1 {
		t.Errorf("run called %d times, want 1", calls)
	}
}

func TestChildPIDExcludesSelf(t *testing.T) {
	c := New(1)
	c.getpid = func() int { return 999 }
	c.run = func(name string, args ...string) string {
		return "999\n123\n"
	}
	if got := c.ChildPID(); got != "123" {
		t.Errorf("ChildPID = %q, want 123", got)
	}
}

func TestChildPIDHighestNonSelf(t *testing.T) {
	c := New(1)
	c.getpid = func() int { return 999 }
	c.run = func(name string, args ...string) string {
		return "123\n456\n"
	}
	if got := c.ChildPID(); got != "456" {
		t.Errorf("ChildPID = %q, want 456", got)
	}
}

func TestChildPIDNone(t *testing.T) {
	c := New(1)
	c.run = func(name string, args ...string) string { return "\n" }
	if got := c.ChildPID(); got != "" {
		t.Errorf("ChildPID = %q, want empty", got)
	}
}
