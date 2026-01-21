package config

func Default() Config {
	return Config{
		Version:            "0.6.6",
		Nick:               "HyLauncher",
		CurrentGameVersion: 0,
		Branch:             "release",
	}
}
