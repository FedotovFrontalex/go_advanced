package configs

const (
	ErrNoDSN             = "no dsn provided to env. You should add dsn string to continue"
	ErrNoSecret          = "no secret phrase is provided. Authorization may work incorectly"
	ErrNoSessionIdLength = "failed to parse session_id length from env. Use 20 by default"
	ErrNoDomains         = "no domains provided to env"
)
