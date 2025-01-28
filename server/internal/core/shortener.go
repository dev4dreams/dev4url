package core

import (
	"errors"
	"sync"
	"time"
)

const (
	// Characters carefully chosen to avoid ambiguity
	alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	// Bit lengths
	timestampBits = 29 // ~17 years
	workerBits    = 8  // 256 workers
	sequenceBits  = 4  // 16 sequences per second

	maxWorkerID = -1 ^ (-1 << workerBits)
	maxSequence = -1 ^ (-1 << sequenceBits)

	// For bit shifting
	workerShift    = sequenceBits
	timestampShift = sequenceBits + workerBits
)

var (
	// Custom epoch (2024-01-01 00:00:00 UTC)
	epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6

	ErrInvalidWorkerID     = errors.New("worker ID exceeds maximum")
	ErrClockMovedBackwards = errors.New("clock moved backwards")
)

// Generator handles the generation of unique IDs
type Generator struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

// NewGenerator creates a new Generator instance
func NewGenerator(workerID int64) (*Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}

	return &Generator{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

// NextID generates a new unique ID
func (g *Generator) NextID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6

	if timestamp < g.timestamp {
		return 0, ErrClockMovedBackwards
	}

	if timestamp == g.timestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// Sequence exhausted, wait for next millisecond
			for timestamp <= g.timestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0
	}

	g.timestamp = timestamp

	id := ((timestamp - epoch) << timestampShift) |
		(g.workerID << workerShift) |
		g.sequence

	return id, nil
}

// GenerateShortURL generates a 7-character short URL
func (g *Generator) GenerateShortURL() (string, error) {
	id, err := g.NextID()
	if err != nil {
		return "", err
	}
	return encodeToBase58(uint64(id)), nil
}

// Validate if a received short URL is legitimate
func (g *Generator) IsValidShortURL(shortURL string) bool {
	_, err := decodeFromBase58(shortURL)
	return err == nil
}

// encodeToBase58 converts a number to base58 string
func encodeToBase58(num uint64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	// Pre-allocate slice with capacity of 7 for efficiency
	chars := make([]byte, 0, 7)
	base := uint64(len(alphabet))

	for num > 0 {
		chars = append([]byte{alphabet[num%base]}, chars...)
		num = num / base
	}

	// Pad to ensure 7 characters
	for len(chars) < 7 {
		chars = append([]byte{alphabet[0]}, chars...)
	}

	// Ensure we don't exceed 7 characters
	if len(chars) > 7 {
		chars = chars[len(chars)-7:]
	}

	return string(chars)
}

// decodeFromBase58 converts a base58 string back to number
func decodeFromBase58(encoded string) (uint64, error) {
	var num uint64
	base := uint64(len(alphabet))

	for _, char := range encoded {
		pos := -1
		for i, c := range alphabet {
			if c == char {
				pos = i
				break
			}
		}
		if pos == -1 {
			return 0, errors.New("invalid character in encoded string")
		}
		num = num*base + uint64(pos)
	}
	return num, nil
}
