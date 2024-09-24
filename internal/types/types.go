package types

// Tried using anonymous structs, but declaring them requires restating their types
type StatusResponse struct {
	Version            SRVersion `json:"version"`
	Players            SRPlayers `json:"players"`
	Description        SRDesc    `json:"description"`
	Favicon            string    `json:"favicon"`
	EnforcesSecureChat bool      `json:"enforcesSecureChat"`
}

type SRVersion struct {
	Name     string `json:"name"`
	Protocol int16  `json:"protocol"`
}

type SRPlayers struct {
	Max    int8 `json:"max"`
	Online int8 `json:"online"`
	Sample []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"sample"`
}

type SRDesc struct {
	Text string `json:"text"`
}
