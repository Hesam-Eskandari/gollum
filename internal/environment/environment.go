package environment

import (
	"flag"
	"fmt"
)

type Environment string

const (
	Prod     Environment = "prod"
	Internal Environment = "internal"
	Develop  Environment = "develop"
)

var defaultEnv Environment = Develop
var environment Environment

func SetEnvironment(env Environment) error {
	if environment != "" {
		return fmt.Errorf("environment is already set to %v", environment)
	}
	if env != Prod && env != Develop && env != Internal {
		environment = defaultEnv
		return fmt.Errorf("invalid environment \"%v\". set to default environment \"%v\"", env, defaultEnv)
	}
	environment = env
	return nil
}

func fromString(env string) (Environment, error) {
	switch env {
	case string(Prod):
		return Prod, nil
	case string(Internal):
		return Internal, nil
	case string(Develop):
		return Develop, nil
	}
	return defaultEnv, fmt.Errorf("invalid environment \"%v\", switching to \"%v\"", env, Develop)
}

func AutoSetEnvironment() error {
	envStr := flag.String("env", string(defaultEnv), "application environment")
	flag.Parse()
	env, err := fromString(*envStr)
	if err != nil {
		return err
	}
	return SetEnvironment(env)
}

func GetEnvironment() Environment {
	return environment
}

func IsDevelop() bool {
	return environment == Develop
}
