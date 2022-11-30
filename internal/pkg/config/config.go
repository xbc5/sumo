package config

type Fetch struct {
	Threads int
}

type Server struct {
	Address string
}

type Config struct {
	Fetch  Fetch
	Server Server
}

func GetConfig() Config {
	return Config{
		Fetch:  Fetch{Threads: 5},
		Server: Server{Address: ":8080"},
	}
}
