package environment

import (
	"flag"
	"fmt"
)

type Environment string

const (
	Prod Environment = "Prod"
	Dev  Environment = "Dev"
)

var environment Environment

func SetEnvironment(env Environment) error {
	if environment != "" {
		return fmt.Errorf("environment is already set to %v", environment)
	}
	if env != Prod && env != Dev {
		return fmt.Errorf("invalid environment %v. choose between %v and %v", env, Dev, Prod)
	}
	environment = env
	return nil
}

func AutoSetEnvironment() error {
	isProdEnvPtr := flag.Bool("prod", false, "is environment prod")
	flag.Parse()
	env := Dev
	if *isProdEnvPtr {
		env = Prod
	}
	return SetEnvironment(env)
}

func GetEnvironment() Environment {
	return environment
}

func IsDevelop() bool {
	return environment == Dev
}
