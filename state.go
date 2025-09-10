package main

import (
	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/config"
	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
