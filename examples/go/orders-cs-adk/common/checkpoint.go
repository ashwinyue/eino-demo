package common

import "context"

type MemCheckPointStore struct {
    buf map[string][]byte
}

func NewMemCheckPointStore() *MemCheckPointStore {
    return &MemCheckPointStore{buf: make(map[string][]byte)}
}

func (m *MemCheckPointStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
    v, ok := m.buf[key]
    return v, ok, nil
}

func (m *MemCheckPointStore) Set(ctx context.Context, key string, value []byte) error {
    m.buf[key] = value
    return nil
}

