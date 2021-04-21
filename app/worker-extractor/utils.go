package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func readYamlFile(filePath string, frame interface{}) error {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to setup configuration from file: %v\n", err)
		return err
	}
	if err = yaml.Unmarshal(raw, frame); err != nil {
		log.Fatalf("failed to setup configuration from file: %v\n", err)
		return err
	}

	return nil
}

func readEnvVariable(env, def string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("defaulting env %s to %s", env, def)
		return def
	}
	return v
}
