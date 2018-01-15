package stringutil_test

import (
	"container/list"
	"github.com/ServiceComb/go-chassis/core/util/string"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var s = strings.Repeat("a", 1024)

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
