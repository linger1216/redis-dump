package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDumpConfig(t *testing.T) {
	dump := &DumpConfig{}
	dump.Common = &commonConfig{
		Parallel: 12,
		Batch:    5000,
	}

	dump.Source = &sourceConfig{
		File: []*sourceFileConfig{
			{
				Filename:   "bk.resp",
				Batch:      dump.Common.Batch,
				ReaderSize: 4096 * 32,
			},
		},
		Single: []*sourceSingleRedisConfig{
			{
				Network:  "tcp",
				Url:      "localhost:6379",
				DBNumber: 0,
				TTL:      true,
				Match:    "*",
			},
		},
		Cluster: []*sourceClusterRedisConfig{
			{
				Url:   []string{"localhost:6379"},
				TTL:   true,
				Match: "*",
			},
		},
	}

	dump.Output = &outputConfig{
		File: &outputFileConfig{
			Flag:      "trunc",
			Filename:  "test.resp",
			WriteSize: 4096 * 32,
		},
		Single: &outputSingleRedisConfig{
			Network:  "tcp",
			Url:      "localhost:6379",
			DBNumber: 0,
		},
		Cluster: &outputClusterRedisConfig{
			Url: []string{"localhost:7379"},
		},
	}

	buf, _ := json.Marshal(dump)
	fmt.Println(string(buf))
}
