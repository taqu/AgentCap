package common

import (
	"crypto/sha256"
	"encoding/hex"
)

// MapAgentSession converts an external agent session ID to a safe local AgentCap session ID.
// The agent session ID may contain sensitive tokens or be too long; we hash it.
// Returns a 6-hex string matching the AgentCap session ID format.
func MapAgentSession(agentSessionID string) string {
	if agentSessionID == "" {
		return ""
	}
	h := sha256.Sum256([]byte(agentSessionID))
	return hex.EncodeToString(h[:])[:6]
}
