package config

func defaultConfig(): Config {
    return Config{
        Port: 8080,
        LogLevel: "DEBUG",
        Stelselnode: Stelselnode{
            Cert: "Blabla",
            Key: "Blabla"
        }
    }
}
