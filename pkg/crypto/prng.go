package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateSortSeed generates a cryptographic seed for the raffle draw.
func GenerateSortSeed(raffleID string) (seedInput string, commitHash string, err error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return "", "", fmt.Errorf("generate nonce: %w", err)
	}

	timestamp := time.Now().UnixNano()
	seedInput = fmt.Sprintf("%s-%x-%d", raffleID, nonce, timestamp)
	commitHash = ComputeCommitHash(seedInput)

	return seedInput, commitHash, nil
}

// ComputeCommitHash computes SHA-256 of the seed input.
// This is published before the draw and can be verified by anyone.
func ComputeCommitHash(seedInput string) string {
	h := sha256.Sum256([]byte(seedInput))
	return hex.EncodeToString(h[:])
}

// PRNG is a deterministic pseudo-random number generator seeded with SHA-256.
type PRNG struct {
	state []byte
}

// NewPRNG creates a new PRNG from a seed string.
func NewPRNG(seed string) *PRNG {
	h := sha256.Sum256([]byte(seed))
	return &PRNG{state: h[:]}
}

// Intn returns a deterministic random integer in [0, n).
// Uses the state bytes cyclically with SHA-256 chaining.
func (p *PRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}

	h := sha256.Sum256(p.state)
	p.state = h[:]

	val := 0
	for i := 0; i < len(h) && i < 8; i++ {
		val = (val << 8) | int(h[i])
	}
	if val < 0 {
		val = -val
	}

	return val % n
}

// FisherYatesShuffle shuffles a slice of integers deterministically using the PRNG.
func FisherYatesShuffle(seed string, tickets []int) {
	prng := NewPRNG(seed)
	for i := len(tickets) - 1; i > 0; i-- {
		j := prng.Intn(i + 1)
		tickets[i], tickets[j] = tickets[j], tickets[i]
	}
}
