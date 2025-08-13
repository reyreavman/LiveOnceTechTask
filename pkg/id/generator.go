package id

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"
)

type Generator struct {
	lastTime  int64
	lastCount uint64
	mu        sync.Mutex
}

func NewGenerator() *Generator {
	return &Generator{
		lastTime:  0,
		lastCount: 0,
	}
}

func (g *Generator) Generate() string {
    buf := make([]byte, 16)
    now := time.Now().UTC().UnixNano()

    binary.BigEndian.PutUint64(buf[:8], uint64(now))

    if _, err := rand.Read(buf[8:]); err == nil {
        return hex.EncodeToString(buf)
    }

    g.mu.Lock()
    defer g.mu.Unlock()
    
    if now == g.lastTime {
        g.lastCount++
    } else {
        g.lastTime = now
        g.lastCount = 0
    }
    
    binary.BigEndian.PutUint64(buf[8:], g.lastCount)
    return hex.EncodeToString(buf)
}
