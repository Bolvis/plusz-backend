package env

import "os"

func Load() {
	env := os.Getenv("ENV")

	switch env {
	case "":
		panic("environment variable not set")
	case "dev":
		if err := os.Setenv("DB_HOST", "localhost"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_PORT", "5432"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_USERNAME", "postgres"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_PASSWORD", "example"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_NAME", "plusz_db"); err != nil {
			panic(err)
		}
		if err := os.Setenv("SSL_MODE", "disable"); err != nil {
			panic(err)
		}
	case "prod":
		if err := os.Setenv("DB_HOST", "ec2-54-195-190-73.eu-west-1.compute.amazonaws.com"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_PORT", "5432"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_USERNAME", "ufqgndot5qvr7m"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_PASSWORD", "ufqgndot5qvr7m"); err != nil {
			panic(err)
		}
		if err := os.Setenv("DB_NAME", "ufqgndot5qvr7m"); err != nil {
			panic(err)
		}
		if err := os.Setenv("SSL_MODE", "require"); err != nil {
			panic(err)
		}
	}
}
