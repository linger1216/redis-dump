package core

import "runtime"

type commonConfig struct {
	Parallel int `json:"parallel"`
	Batch    int `json:"batch"`
}

func (c *commonConfig) Default() {
	if c.Parallel == 0 {
		c.Parallel = runtime.NumCPU()
	}
	if c.Batch == 0 {
		c.Batch = 5000
	}
}

type sourceConfig struct {
	File    sourceFileConfig         `json:"file"`
	Single  sourceSingleRedisConfig  `json:"single"`
	Cluster sourceClusterRedisConfig `json:"cluster"`
}

type outputConfig struct {
	File    outputFileConfig         `json:"file"`
	Single  outputSingleRedisConfig  `json:"single"`
	Cluster outputClusterRedisConfig `json:"cluster"`
}

type DumpConfig struct {
	Common  commonConfig   `json:"common"`
	Sources []sourceConfig `json:"sources"`
	Outputs []outputConfig `json:"outputs"`
}
