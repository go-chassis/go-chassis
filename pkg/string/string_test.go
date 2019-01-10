package stringutil_test

import (
	"container/list"
	"github.com/go-chassis/go-chassis/pkg/string"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var s = strings.Repeat("a", 1024)

func TestClearByteMemory(t *testing.T) {
	b := []byte("aaa")
	assert.Equal(t, "aaa", string(b))
	stringutil.ClearByteMemory(b)
	assert.NotEqual(t, "aaa", string(b))
}
func TestMinInt(t *testing.T) {

	a := stringutil.MinInt(1, 2)
	assert.Equal(t, 1, a)

	a = stringutil.MinInt(2, 1)
	assert.Equal(t, 1, a)
}
func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := []byte(s)
		_ = string(b)
	}
}

func BenchmarkTestBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := []byte(s)
		_ = stringutil.Bytes2str(b)
	}
}
func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte(s)

	}
}

func BenchmarkTestBlock2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringutil.Str2bytes("asd")
	}
}

const capacity = 2

func array() [capacity]int {
	var d [capacity]int

	for i := 0; i < len(d); i++ {
		d[i] = 1
	}

	return d
}

func slice() []int {
	d := make([]int, capacity)

	for i := 0; i < len(d); i++ {
		d[i] = 1
	}

	return d
}

var l *list.List

func init() {
	l = list.New()
	for i := 0; i < capacity; i++ {
		l.PushBack(i)
	}
}
func List() {

	for e := l.Front(); e != nil; e = e.Next() {
		_ = e.Value.(int)
	}

}
func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = array()
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = slice()
	}

}
func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		List()
	}

}

func TestStringFunc(t *testing.T) {

	var strArr = []string{"abc", "def"}

	b := stringutil.StringInSlice("abc", strArr)
	assert.Equal(t, true, b)

	b = stringutil.StringInSlice("wer", strArr)
	assert.Equal(t, false, b)

	by := stringutil.Str2bytes("abc")
	str := stringutil.Bytes2str(by)

	assert.Equal(t, str, string(by))
}

func TestSplitToTwo(t *testing.T) {
	s1 := "aa"
	s2 := "bb"
	sub := "::"
	s := s1 + sub + s2
	p1, p2 := stringutil.SplitToTwo(s, sub)
	assert.Equal(t, s1, p1)
	assert.Equal(t, s2, p2)

	p1, p2 = stringutil.SplitToTwo(s, "/")
	assert.Empty(t, p1)
	assert.Equal(t, s, p2)
}

func TestSplitFirstSep(t *testing.T) {
	s1 := "aa"
	s2 := "bb"
	sub := "::"
	s := s1 + sub + s2
	p1 := stringutil.SplitFirstSep(s, sub)
	assert.Equal(t, s1, p1)

	p1 = stringutil.SplitFirstSep(s, "/")
	assert.Empty(t, p1)
}

func SplitToTwoByStringsSplit(s, sep string) (string, string) {
	r := strings.Split(s, sep)
	return r[0], r[1]
}

func SplitFirstSepByStringsSplit(s, sep string) string {
	return strings.Split(s, sep)[0]
}

var testURL = "http://127.0.0.1"

func Benchmark_SplitToTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = stringutil.SplitToTwo(testURL, "://")
	}
}

func Benchmark_SplitToTwoByStringsSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = SplitToTwoByStringsSplit(testURL, "://")
	}
}

func Benchmark_SplitFirstSep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stringutil.SplitFirstSep(testURL, "://")
	}
}

func Benchmark_SplitFirstSepByStringsSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SplitFirstSepByStringsSplit(testURL, "://")
	}
}
