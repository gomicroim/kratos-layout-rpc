package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func MustLoad(flagconf string) *Bootstrap {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
