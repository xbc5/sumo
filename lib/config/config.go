package config

type Fetch struct {
	Threads int
}

type Config struct {
	Fetch Fetch
}

func GetConfig() Config {
	return Config{
		Fetch: Fetch{Threads: 5},
	}
}
