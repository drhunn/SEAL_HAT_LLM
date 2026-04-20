package slotpacket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var familyIDs = map[string]int{
	"constitutional": 0,
	"tools":          1,
	"skills":         2,
	"prompt":         3,
	"heartbeat":      4,
	"memory":         5,
	"dreams":         6,
	"postmortem":     7,
}

var authorityIDs = map[string]int{
	"parent":     0,
	"harness":    1,
	"specialist": 2,
}

var familyBySlotName = map[string]string{
	"IDENTITY.md":   "constitutional",
	"SOUL.md":       "constitutional",
	"AGENTS.md":     "constitutional",
	"TOOLS.md":      "tools",
	"SKILLS.md":     "skills",
	"PROMPT.md":     "prompt",
	"HEARTBEAT.md":  "heartbeat",
	"MEMORY.md":     "memory",
	"DREAMS.md":     "dreams",
	"POSTMORTEM.md": "postmortem",
}

func BuildFromUnknown(specialistID, mode string, parentFrozen bool, raw any) Packet {
	records := extractRecords(raw)
	slotTexts := make([]string, 0, len(records))
	slotFamilyIDs := make([]int, 0, len(records))
	slotAuthorityIDs := make([]int, 0, len(records))
	slotEnabledMask := make([]int, 0, len(records))

	for _, record := range records {
		slotTexts = append(slotTexts, record.Text)
		slotFamilyIDs = append(slotFamilyIDs, familyIDs[record.Family])
		slotAuthorityIDs = append(slotAuthorityIDs, authorityIDs[record.Authority])
		slotEnabledMask = append(slotEnabledMask, 1)
	}

	packet := Packet{
		Version:          1,
		SpecialistID:     specialistID,
		Mode:             mode,
		ParentFrozen:     parentFrozen,
		SlotTexts:        slotTexts,
		SlotFamilyIDs:    slotFamilyIDs,
		SlotAuthorityIDs: slotAuthorityIDs,
		SlotEnabledMask:  slotEnabledMask,
		ScalarFlags: []int{
			1,
			0,
			1,
			1,
		},
		Records: records,
	}
	packet.VersionHash = versionHash(packet)
	return packet
}

func extractRecords(raw any) []Record {
	v := reflect.ValueOf(raw)
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	records := make([]Record, 0)
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
		})
		for _, key := range keys {
			record, ok := recordFromPair(fmt.Sprint(key.Interface()), v.MapIndex(key).Interface())
			if ok {
				records = append(records, record)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			record, ok := recordFromUnknown(v.Index(i).Interface())
			if ok {
				records = append(records, record)
			}
		}
	default:
		record, ok := recordFromUnknown(v.Interface())
		if ok {
			records = append(records, record)
		}
	}
	return records
}

func recordFromPair(name string, value any) (Record, bool) {
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "" {
		return Record{}, false
	}
	return buildRecord(normalizeSlotName(name), text), true
}

func recordFromUnknown(raw any) (Record, bool) {
	v := reflect.ValueOf(raw)
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return Record{}, false
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return Record{}, false
	}
	if v.Kind() != reflect.Struct {
		return Record{}, false
	}

	name := firstStringField(v, []string{"SlotName", "Name", "Filename", "FileName", "Path"})
	text := firstStringField(v, []string{"Text", "Content", "Body", "Value"})
	if strings.TrimSpace(name) == "" || strings.TrimSpace(text) == "" {
		return Record{}, false
	}
	return buildRecord(normalizeSlotName(name), text), true
}

func buildRecord(slotName, text string) Record {
	family := familyBySlotName[slotName]
	if family == "" {
		family = "prompt"
	}
	authority := "harness"
	hotReloadable := true
	if family == "constitutional" {
		authority = "parent"
		hotReloadable = false
	}
	return Record{
		SlotName:      slotName,
		Family:        family,
		Authority:     authority,
		Text:          text,
		Summary:       summaryFromText(text),
		Rules:         rulesFromText(text),
		HotReloadable: hotReloadable,
	}
}

func normalizeSlotName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\\", "/")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return strings.ToUpper(name)
}

func firstStringField(v reflect.Value, names []string) string {
	for _, name := range names {
		f := v.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}
	return ""
}

func summaryFromText(text string) string {
	lines := make([]string, 0, 2)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) == 2 {
			break
		}
	}
	summary := strings.Join(lines, " ")
	if len(summary) > 240 {
		return summary[:240]
	}
	return summary
}

func rulesFromText(text string) []string {
	rules := make([]string, 0, 16)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			rules = append(rules, strings.TrimSpace(strings.TrimPrefix(line, "-")))
		}
		if len(rules) == 16 {
			break
		}
	}
	return rules
}

func versionHash(packet Packet) string {
	h := sha256.New()
	for _, record := range packet.Records {
		_, _ = h.Write([]byte(record.SlotName))
		_, _ = h.Write([]byte{0})
		_, _ = h.Write([]byte(record.Text))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
