package slotpacket

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func WriteJSON(path string, packet Packet) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(packet, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
