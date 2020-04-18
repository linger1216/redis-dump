package core

import "runtime"

type commonConfig struct {
	Parallel int
	Batch    int `default:"5000"`
}

func (c *commonConfig) Default() {
	if c.Parallel == 0 {
		c.Parallel = runtime.NumCPU()
	}
	if c.Batch == 0 {
		c.Batch = 5000
	}
}

type redisConfig struct {
	Network  string
	Url      []string
	Password string
	DBNumber int
	TTL      bool
	Match    string
}

func (c *redisConfig) Default() {
	if len(c.Network) == 0 {
		c.Network = "tcp"
	}

	if len(c.Match) == 0 {
		c.Match = "*"
	}
}

// --------------------------------------------------------------------------------

type destFileConfig struct {
	Out      string
	FileName string
}

func (c *destFileConfig) Default() {
	if len(c.Out) == 0 {
		c.Out = "resp"
	}
}

type destRedisConfig struct {
	Network  string
	Url      []string
	Password string
}

func (c *destRedisConfig) Default() {
	if len(c.Network) == 0 {
		c.Network = "tcp"
	}
}

type destConfig struct {
	File  destFileConfig
	Redis destRedisConfig
}

func (c *destConfig) Default() {
	c.Redis.Default()
	c.File.Default()
}

type DumpConfig struct {
	Common commonConfig
	Src    redisConfig
	Dest   destConfig
}

func (c *DumpConfig) Default() {
	c.Common.Default()
	c.Src.Default()
	c.Dest.Default()
}
