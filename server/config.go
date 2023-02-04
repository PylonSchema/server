package server

type conf struct {
	Database *databaseInfo
	Sentry   *sentryInfo
	Secret   *secret
	Oauth    map[string]oauth2Info
}

type databaseInfo struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	Port     string `toml:"port"`
}

type sentryInfo struct {
	Dsn string
}

type secret struct {
	Session string
	Jwtkey  string
}

type oauth2Info struct {
	Id       string
	Secret   string
	Redirect string
}
