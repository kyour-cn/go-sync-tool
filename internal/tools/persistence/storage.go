package persistence

import (
    "bytes"
    "encoding/gob"
    "errors"
    "fmt"
    "io/fs"
    "os"
)

// Storage 提供类型安全的数据存储功能
type Storage[T any] struct{}

// NewStorage 创建新的存储实例
func NewStorage[T any]() *Storage[T] {
    return &Storage[T]{}
}

// Save 将数据保存到文件
func (s *Storage[T]) Save(data T, filename string) error {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    if err := enc.Encode(data); err != nil {
        return fmt.Errorf("encode failed: %w", err)
    }

    return os.WriteFile(filename, buf.Bytes(), 0o644)
}

// Load 从文件加载数据
func (s *Storage[T]) Load(filename string) (T, error) {
    var zero T

    data, err := os.ReadFile(filename)
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            return zero, fmt.Errorf("file %q does not exist: %w", filename, err)
        }
        return zero, fmt.Errorf("read file failed: %w", err)
    }

    dec := gob.NewDecoder(bytes.NewReader(data))
    err = dec.Decode(&zero)
    if err != nil {
        return zero, fmt.Errorf("decode failed: %w", err)
    }

    return zero, nil
}
