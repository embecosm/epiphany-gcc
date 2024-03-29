// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"testing"
)

type deflateTest struct {
	in    []byte
	level int
	out   []byte
}

type deflateInflateTest struct {
	in []byte
}

type reverseBitsTest struct {
	in       uint16
	bitCount uint8
	out      uint16
}

var deflateTests = []*deflateTest{
	{[]byte{}, 0, []byte{1, 0, 0, 255, 255}},
	{[]byte{0x11}, -1, []byte{18, 4, 4, 0, 0, 255, 255}},
	{[]byte{0x11}, DefaultCompression, []byte{18, 4, 4, 0, 0, 255, 255}},
	{[]byte{0x11}, 4, []byte{18, 4, 4, 0, 0, 255, 255}},

	{[]byte{0x11}, 0, []byte{0, 1, 0, 254, 255, 17, 1, 0, 0, 255, 255}},
	{[]byte{0x11, 0x12}, 0, []byte{0, 2, 0, 253, 255, 17, 18, 1, 0, 0, 255, 255}},
	{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 0,
		[]byte{0, 8, 0, 247, 255, 17, 17, 17, 17, 17, 17, 17, 17, 1, 0, 0, 255, 255},
	},
	{[]byte{}, 1, []byte{1, 0, 0, 255, 255}},
	{[]byte{0x11}, 1, []byte{18, 4, 4, 0, 0, 255, 255}},
	{[]byte{0x11, 0x12}, 1, []byte{18, 20, 2, 4, 0, 0, 255, 255}},
	{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 1, []byte{18, 132, 2, 64, 0, 0, 0, 255, 255}},
	{[]byte{}, 9, []byte{1, 0, 0, 255, 255}},
	{[]byte{0x11}, 9, []byte{18, 4, 4, 0, 0, 255, 255}},
	{[]byte{0x11, 0x12}, 9, []byte{18, 20, 2, 4, 0, 0, 255, 255}},
	{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 9, []byte{18, 132, 2, 64, 0, 0, 0, 255, 255}},
}

var deflateInflateTests = []*deflateInflateTest{
	{[]byte{}},
	{[]byte{0x11}},
	{[]byte{0x11, 0x12}},
	{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
	{[]byte{0x11, 0x10, 0x13, 0x41, 0x21, 0x21, 0x41, 0x13, 0x87, 0x78, 0x13}},
	{largeDataChunk()},
}

var reverseBitsTests = []*reverseBitsTest{
	{1, 1, 1},
	{1, 2, 2},
	{1, 3, 4},
	{1, 4, 8},
	{1, 5, 16},
	{17, 5, 17},
	{257, 9, 257},
	{29, 5, 23},
}

func largeDataChunk() []byte {
	result := make([]byte, 100000)
	for i := range result {
		result[i] = byte(i * i & 0xFF)
	}
	return result
}

func TestDeflate(t *testing.T) {
	for _, h := range deflateTests {
		var buf bytes.Buffer
		w := NewWriter(&buf, h.level)
		w.Write(h.in)
		w.Close()
		if !bytes.Equal(buf.Bytes(), h.out) {
			t.Errorf("Deflate(%d, %x) = %x, want %x", h.level, h.in, buf.Bytes(), h.out)
		}
	}
}

type syncBuffer struct {
	buf    bytes.Buffer
	mu     sync.RWMutex
	closed bool
	ready  chan bool
}

func newSyncBuffer() *syncBuffer {
	return &syncBuffer{ready: make(chan bool, 1)}
}

func (b *syncBuffer) Read(p []byte) (n int, err error) {
	for {
		b.mu.RLock()
		n, err = b.buf.Read(p)
		b.mu.RUnlock()
		if n > 0 || b.closed {
			return
		}
		<-b.ready
	}
	panic("unreachable")
}

func (b *syncBuffer) signal() {
	select {
	case b.ready <- true:
	default:
	}
}

func (b *syncBuffer) Write(p []byte) (n int, err error) {
	n, err = b.buf.Write(p)
	b.signal()
	return
}

func (b *syncBuffer) WriteMode() {
	b.mu.Lock()
}

func (b *syncBuffer) ReadMode() {
	b.mu.Unlock()
	b.signal()
}

func (b *syncBuffer) Close() error {
	b.closed = true
	b.signal()
	return nil
}

func testSync(t *testing.T, level int, input []byte, name string) {
	if len(input) == 0 {
		return
	}

	t.Logf("--testSync %d, %d, %s", level, len(input), name)
	buf := newSyncBuffer()
	buf1 := new(bytes.Buffer)
	buf.WriteMode()
	w := NewWriter(io.MultiWriter(buf, buf1), level)
	r := NewReader(buf)

	// Write half the input and read back.
	for i := 0; i < 2; i++ {
		var lo, hi int
		if i == 0 {
			lo, hi = 0, (len(input)+1)/2
		} else {
			lo, hi = (len(input)+1)/2, len(input)
		}
		t.Logf("#%d: write %d-%d", i, lo, hi)
		if _, err := w.Write(input[lo:hi]); err != nil {
			t.Errorf("testSync: write: %v", err)
			return
		}
		if i == 0 {
			if err := w.Flush(); err != nil {
				t.Errorf("testSync: flush: %v", err)
				return
			}
		} else {
			if err := w.Close(); err != nil {
				t.Errorf("testSync: close: %v", err)
			}
		}
		buf.ReadMode()
		out := make([]byte, hi-lo+1)
		m, err := io.ReadAtLeast(r, out, hi-lo)
		t.Logf("#%d: read %d", i, m)
		if m != hi-lo || err != nil {
			t.Errorf("testSync/%d (%d, %d, %s): read %d: %d, %v (%d left)", i, level, len(input), name, hi-lo, m, err, buf.buf.Len())
			return
		}
		if !bytes.Equal(input[lo:hi], out[:hi-lo]) {
			t.Errorf("testSync/%d: read wrong bytes: %x vs %x", i, input[lo:hi], out[:hi-lo])
			return
		}
		// This test originally checked that after reading
		// the first half of the input, there was nothing left
		// in the read buffer (buf.buf.Len() != 0) but that is
		// not necessarily the case: the write Flush may emit
		// some extra framing bits that are not necessary
		// to process to obtain the first half of the uncompressed
		// data.  The test ran correctly most of the time, because
		// the background goroutine had usually read even
		// those extra bits by now, but it's not a useful thing to
		// check.
		buf.WriteMode()
	}
	buf.ReadMode()
	out := make([]byte, 10)
	if n, err := r.Read(out); n > 0 || err != io.EOF {
		t.Errorf("testSync (%d, %d, %s): final Read: %d, %v (hex: %x)", level, len(input), name, n, err, out[0:n])
	}
	if buf.buf.Len() != 0 {
		t.Errorf("testSync (%d, %d, %s): extra data at end", level, len(input), name)
	}
	r.Close()

	// stream should work for ordinary reader too
	r = NewReader(buf1)
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("testSync: read: %s", err)
		return
	}
	r.Close()
	if !bytes.Equal(input, out) {
		t.Errorf("testSync: decompress(compress(data)) != data: level=%d input=%s", level, name)
	}
}

func testToFromWithLevel(t *testing.T, level int, input []byte, name string) error {
	buffer := bytes.NewBuffer(nil)
	w := NewWriter(buffer, level)
	w.Write(input)
	w.Close()
	r := NewReader(buffer)
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("read: %s", err)
		return err
	}
	r.Close()
	if !bytes.Equal(input, out) {
		t.Errorf("decompress(compress(data)) != data: level=%d input=%s", level, name)
	}

	testSync(t, level, input, name)
	return nil
}

func testToFrom(t *testing.T, input []byte, name string) {
	for i := 0; i < 10; i++ {
		testToFromWithLevel(t, i, input, name)
	}
}

func TestDeflateInflate(t *testing.T) {
	for i, h := range deflateInflateTests {
		testToFrom(t, h.in, fmt.Sprintf("#%d", i))
	}
}

func TestReverseBits(t *testing.T) {
	for _, h := range reverseBitsTests {
		if v := reverseBits(h.in, h.bitCount); v != h.out {
			t.Errorf("reverseBits(%v,%v) = %v, want %v",
				h.in, h.bitCount, v, h.out)
		}
	}
}

func TestDeflateInflateString(t *testing.T) {
	gold, err := ioutil.ReadFile("../testdata/e.txt")
	if err != nil {
		t.Error(err)
	}
	testToFromWithLevel(t, 1, gold, "2.718281828...")
}

func TestReaderDict(t *testing.T) {
	const (
		dict = "hello world"
		text = "hello again world"
	)
	var b bytes.Buffer
	w := NewWriter(&b, 5)
	w.Write([]byte(dict))
	w.Flush()
	b.Reset()
	w.Write([]byte(text))
	w.Close()

	r := NewReaderDict(&b, []byte(dict))
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello again world" {
		t.Fatalf("read returned %q want %q", string(data), text)
	}
}

func TestWriterDict(t *testing.T) {
	const (
		dict = "hello world"
		text = "hello again world"
	)
	var b bytes.Buffer
	w := NewWriter(&b, 5)
	w.Write([]byte(dict))
	w.Flush()
	b.Reset()
	w.Write([]byte(text))
	w.Close()

	var b1 bytes.Buffer
	w = NewWriterDict(&b1, 5, []byte(dict))
	w.Write([]byte(text))
	w.Close()

	if !bytes.Equal(b1.Bytes(), b.Bytes()) {
		t.Fatalf("writer wrote %q want %q", b1.Bytes(), b.Bytes())
	}
}

// See http://code.google.com/p/go/issues/detail?id=2508
func TestRegression2508(t *testing.T) {
	w := NewWriter(ioutil.Discard, 1)
	buf := make([]byte, 1024)
	for i := 0; i < 131072; i++ {
		if _, err := w.Write(buf); err != nil {
			t.Fatalf("writer failed: %v", err)
		}
	}
	w.Close()
}
