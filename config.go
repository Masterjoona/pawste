package main

type ConfigProgram struct {
	PasswordSalt string
}

var Config ConfigProgram

func InitConfig() {
	// This function should read the configuration from a file or environment
	// variables. For now, we'll just hardcode the values.
	Config = ConfigProgram{
		PasswordSalt: "banana",
	}
}
