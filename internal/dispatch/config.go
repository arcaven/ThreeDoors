package dispatch

// DispatchConfig holds configuration for the dev dispatch pipeline.
type DispatchConfig struct {
	// RequireStory controls whether story files are generated before dispatching.
	// Defaults to false — when false, workers receive raw task descriptions.
	RequireStory bool `yaml:"require_story" json:"require_story"`
}
