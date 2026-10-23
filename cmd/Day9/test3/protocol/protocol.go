package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MaxBodySize 限制一条消息的 Body 最多为 1 MiB。
const MaxBodySize = 1 << 20

// WriteMessage 写出一帧：[4 字节大端 Body 长度][Body]。
// 同一连接的消息应由一个 goroutine 顺序写入。
func WriteMessage(w io.Writer, body []byte) error {
	if len(body) > MaxBodySize {
		return fmt.Errorf("消息过长：%d 字节，上限 %d", len(body), MaxBodySize)
	}

	frame := make([]byte, 4+len(body))
	binary.BigEndian.PutUint32(frame[:4], uint32(len(body)))
	copy(frame[4:], body)

	n, err := w.Write(frame)
	if err != nil {
		return err
	}
	if n != len(frame) {
		return io.ErrShortWrite
	}
	return nil
}

// ReadMessage 读取一帧；成功时返回完整 Body，不包含长度头。
func ReadMessage(r io.Reader) ([]byte, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header[:])
	fmt.Printf("length:%#v", length)
	// 先检查长度，再分配内存。
	if length > MaxBodySize {
		return nil, fmt.Errorf("消息过长：%d 字节，上限 %d", length, MaxBodySize)
	}

	body := make([]byte, int(length))
	if _, err := io.ReadFull(r, body); err != nil {
		// 已收到完整消息头，Body 完全缺失也属于消息截断。
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return body, nil
}
