package cli

// Exit codes for CLI commands. All CLI commands should use these
// constants instead of raw integer literals.
const (
	ExitSuccess        = 0
	ExitGeneralError   = 1
	ExitNotFound       = 2
	ExitValidation     = 3
	ExitProviderError  = 4
	ExitAmbiguousInput = 5
)
