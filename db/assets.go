package db

import "embed"

//go:embed migrations
var assets embed.FS

func Assets() embed.FS {
	return assets
}
