package db

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) func(*Postgres) {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) func(*Postgres) {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) func(*Postgres) {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
