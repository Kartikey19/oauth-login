package config

// Config : App wide configuration container
type Config struct {
	// postgres config
	PostgrestHost    string `env:"POSTGRES_HOST"`
	PostgrestPort    string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPass     string `env:"POSTGRES_PASS"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`

	// Logging flags
	Verbose string `env:"VERBOSE_MODE" envDefault:"false"`

	//Host
	Host string `env:"HOST" envDefault:""`

	// Github Variables
	GithubClientID     string `env:"GITHUB_CLIENT_ID" envDefault:""`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET" envDefault:""`

	// Twitter Variables
	TwitterConsumerKey    string `env:"TWITTER_CONSUMER_KEY" envDefault:""`
	TwitterConsumerSecret string `env:"TWITTER_CONSUMER_SECRET" envDefault:""`

	// facebook variables
	FacebookClientID     string `env:"FACEBOOK_CLIENT_ID" envDefault:""`
	FacebookClientSecret string `env:"FACEBOOK_CLIENT_SECRET" envDefault:""`
}
