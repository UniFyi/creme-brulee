package config

import (
	"fmt"
	"net/url"
)

type PsqlConfig struct {
	PostgresHost         string
	PostgresPort         string
	PostgresUser         string
	PostgresPassword     string
	PostgresDatabaseName string
	PostgresSslMode      string
	PostgresSslRootCert  string
}

func NewPsqlConfig() (*PsqlConfig, error) {

	envPostgresHost, err := GetEnv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}
	envPostgresPort, err := GetEnv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}
	envPostgresUser, err := GetEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}
	envPostgresPassword, err := GetEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}
	envPostgresDatabaseName, err := GetEnv("POSTGRES_DATABASE_NAME")
	if err != nil {
		return nil, err
	}
	envPostgresSslMode, err := GetEnv("POSTGRES_SSL_MODE")
	if err != nil {
		return nil, err
	}
	var envPostgresSslRootCertificate string
	if envPostgresSslMode != "disable" {
		envPostgresSslRootCertificate, err = GetEnv("POSTGRES_SSL_ROOT_CERT")
		if err != nil {
			return nil, err
		}
	}

	return &PsqlConfig{
		PostgresHost:         envPostgresHost,
		PostgresPort:         envPostgresPort,
		PostgresUser:         envPostgresUser,
		PostgresPassword:     envPostgresPassword,
		PostgresDatabaseName: envPostgresDatabaseName,
		PostgresSslMode:      envPostgresSslMode,
		PostgresSslRootCert:  envPostgresSslRootCertificate,
	}, nil
}

func (c *PsqlConfig) GetDataSourcePSQL() *url.URL {
	postgresURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:     fmt.Sprintf("%v:%v", c.PostgresHost, c.PostgresPort),
		Path:     fmt.Sprintf("/%v", c.PostgresDatabaseName),
		RawQuery: "sslmode=disable",
	}

	if c.PostgresSslMode != "disable" {
		postgresURL.RawQuery = fmt.Sprintf("sslmode=%v&sslrootcert=%v", c.PostgresSslMode, c.PostgresSslRootCert)
	}

	return postgresURL
}
