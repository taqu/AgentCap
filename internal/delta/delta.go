// Package delta implements session-aware result comparison and delta generation.
package delta

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/taqu/agentcap/internal/store"
)

// Presentation describes how a result was presented to the agent.
type Presentation string

const (
	PresentationFull      Presentation = "full"
	PresentationUnchanged Presentation = "unchanged"
	PresentationDelta     Presentation = "delta"
)

// Result holds the output of the delta engine.
type Result struct {
	Output       string
	Presentation Presentation
}

// Reducer is optionally implemented by delta engines that support structural comparison.
type Reducer interface {
	Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error)
}

// deltaThreshold: emit delta only when smaller than this fraction of the full capsule.
const deltaThreshold = 0.8

// HashBytes returns the sha256 hex of b.
func HashBytes(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
