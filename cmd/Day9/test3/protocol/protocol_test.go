package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

// 模拟每次底层 Read 最多只能取得一个字节，覆盖 Header 和 Body 拆分。
type oneByteReader struct {
	r io.Reader
}

func (r oneByteReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return r.r.Read(p)
}

func TestMessageBoundaries(t *testing.T) {
	messages := []string{"hello", "你好", "第一行\n第二行", "", "world"}
	var stream bytes.Buffer
	for _, message := range messages {
		if err := WriteMessage(&stream, []byte(message)); err != nil {
			t.Fatal(err)
		}
	}

	// 同一字节流包含多帧，同时分别验证一次可读多字节和每次只读一字节。
	for _, fragmented := range []bool{false, true} {
		t.Run(map[bool]string{false: "combined", true: "fragmented"}[fragmented], func(t *testing.T) {
			var reader io.Reader = bytes.NewReader(stream.Bytes())
			if fragmented {
				reader = oneByteReader{r: reader}
			}
			for _, want := range messages {
				body, err := ReadMessage(reader)
				if err != nil || string(body) != want {
					t.Fatalf("got %q, %v; want %q", body, err, want)
				}
			}
			if _, err := ReadMessage(reader); !errors.Is(err, io.EOF) {
				t.Fatalf("want EOF after all frames, got %v", err)
			}
		})
	}
}

func TestWireFormat(t *testing.T) {
	var stream bytes.Buffer
	if err := WriteMessage(&stream, []byte("你好")); err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0, 0, 0, 6}, []byte("你好")...)
	if !bytes.Equal(stream.Bytes(), want) {
		t.Fatalf("got %x, want %x", stream.Bytes(), want)
	}
}

func TestTruncatedMessage(t *testing.T) {
	for _, data := range [][]byte{
		{0, 0},            // Header 不完整。
		{0, 0, 0, 5},      // Header 完整，Body 缺失。
		{0, 0, 0, 5, 'h'}, // Body 不完整。
	} {
		if _, err := ReadMessage(bytes.NewReader(data)); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("data %x: want unexpected EOF, got %v", data, err)
		}
	}
}

func TestOversizedMessage(t *testing.T) {
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], MaxBodySize+1)
	if _, err := ReadMessage(bytes.NewReader(header[:])); err == nil {
		t.Fatal("oversized incoming message must fail before reading Body")
	}
	var stream bytes.Buffer
	if err := WriteMessage(&stream, make([]byte, MaxBodySize+1)); err == nil {
		t.Fatal("oversized outgoing message must fail")
	}
	if stream.Len() != 0 {
		t.Fatal("oversized outgoing message must not write bytes")
	}
}
