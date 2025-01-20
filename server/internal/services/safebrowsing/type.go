package safebrowsing

// ThreatRequest represents the request body for the Safe Browsing API
type ThreatRequest struct {
	Client     ClientInfo `json:"client"`
	ThreatInfo ThreatInfo `json:"threatInfo"`
}

// ClientInfo contains the client identification information
type ClientInfo struct {
	ClientID      string `json:"clientId"`
	ClientVersion string `json:"clientVersion"`
}

// ThreatInfo contains the threat information for the request
type ThreatInfo struct {
	ThreatTypes      []string      `json:"threatTypes"`
	PlatformTypes    []string      `json:"platformTypes"`
	ThreatEntryTypes []string      `json:"threatEntryTypes"`
	ThreatEntries    []ThreatEntry `json:"threatEntries"`
}

// ThreatEntry represents a URL to be checked
type ThreatEntry struct {
	URL string `json:"url"`
}

// ThreatResponse represents the response from the Safe Browsing API
type ThreatResponse struct {
	Matches []ThreatMatch `json:"matches"`
}

// ThreatMatch represents a match found in the Safe Browsing database
type ThreatMatch struct {
	ThreatType          string              `json:"threatType"`
	PlatformType        string              `json:"platformType"`
	ThreatEntryType     string              `json:"threatEntryType"`
	Threat              ThreatEntry         `json:"threat"`
	ThreatEntryMetadata ThreatEntryMetadata `json:"threatEntryMetadata"`
	CacheDuration       string              `json:"cacheDuration"`
}

// ThreatEntryMetadata contains additional information about the threat
type ThreatEntryMetadata struct {
	Entries []MetadataEntry `json:"entries"`
}

// MetadataEntry represents a key-value pair in the threat metadata
type MetadataEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
