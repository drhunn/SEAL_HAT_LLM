package slotpacket

type Record struct {
	SlotName      string   `json:"slot_name"`
	Family        string   `json:"family"`
	Authority     string   `json:"authority"`
	Text          string   `json:"text"`
	Summary       string   `json:"summary"`
	Rules         []string `json:"rules,omitempty"`
	HotReloadable bool     `json:"hot_reloadable"`
}

type Packet struct {
	Version          int      `json:"version"`
	SpecialistID     string   `json:"specialist_id"`
	Mode             string   `json:"mode"`
	ParentFrozen     bool     `json:"parent_frozen"`
	VersionHash      string   `json:"version_hash"`
	SlotTexts        []string `json:"slot_texts"`
	SlotFamilyIDs    []int    `json:"slot_family_ids"`
	SlotAuthorityIDs []int    `json:"slot_authority_ids"`
	SlotEnabledMask  []int    `json:"slot_enabled_mask"`
	ScalarFlags      []int    `json:"scalar_flags"`
	Records          []Record `json:"records"`
}
