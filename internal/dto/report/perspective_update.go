package report

type Score struct {
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

type SpanScore struct {
	Begin int   `json:"begin"`
	End   int   `json:"end"`
	Score Score `json:"score"`
}

type Attribute struct {
	SpanScores   []SpanScore `json:"spanScores"`
	SummaryScore Score       `json:"summaryScore"`
}

type AttributeScores struct {
	Toxicity       Attribute `json:"TOXICITY"`
	SevereToxicity Attribute `json:"SEVERE_TOXICITY"`
	IdentityAttack Attribute `json:"IDENTITY_ATTACK"`
	Insult         Attribute `json:"INSULT"`
	Profanity      Attribute `json:"PROFANITY"`
	Threat         Attribute `json:"THREAT"`
}

type PerspectiveAPIResponse struct {
	AttributeScores   AttributeScores `json:"attributeScores"`
	Languages         []string        `json:"languages"`
	DetectedLanguages []string        `json:"detectedLanguages"`
}
