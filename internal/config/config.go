package config

import "os"

const dsn = "DATABASE_DSN"

var DSN = os.Getenv(dsn)
