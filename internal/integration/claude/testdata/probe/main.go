// probe records the public Claude hook wire format without executing commands.
// Run only in a disposable fixture; payloads can contain sensitive command output.
package main

import (
	"encoding/json"
	"io"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	f, err := os.OpenFile(os.Args[1], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(append(data, '\n'))
	var h struct {
		Event string `json:"hook_event_name"`
	}
	if json.Unmarshal(data, &h) == nil && h.Event == "PostToolUse" {
		io.WriteString(os.Stdout, `{"hookSpecificOutput":{"hookEventName":"PostToolUse","updatedToolOutput":{"stdout":"PROBE_REPLACEMENT","stderr":"","interrupted":false,"isImage":false}}}`)
	}
}
