package config

type (
	// Config configures this server
	Config struct {
		APIUrl       string `toml:"APIUrl"`
		APIToken     string `toml:"APIToken"`
		Host         string
		Secret       string
		ClientID     string
		ClientSecret string
		Permissions  []Permission
		Debug        bool
	}

	// Permission represents the permission levels for a user
	Permission struct {
		ID    string
		Level int
	}
)
