package handler

import "HelpStudent/core/store/rds"

var (
	cache *rds.Redis
)

func Init(c *rds.Redis) {
	cache = c
}
