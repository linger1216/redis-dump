package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDumpConfig(t *testing.T) {
	dump := DumpConfig{}
	dump.Common.Batch = 6000
	dump.Common.Parallel = 12

	/*
		sourceFileConfig{
				Filename:   "bk.resp",
				Batch:      dump.Common.Batch,
				ReaderSize: 4096 * 32,
			}

		sourceSingleRedisConfig{
				Network:  "tcp",
				Url:      "localhost:6379",
				DBNumber: 0,
				TTL:      true,
				Match:    "*",
			}

		sourceClusterRedisConfig{
				Url:   []string{"localhost:6379"},
				TTL:   true,
				Match: "*",
			}



				outputFileConfig{
					Flag:      "trunc",
					Filename:  "test.resp",
					WriteSize: 4096 * 32,
				}

				outputSingleRedisConfig{
					Network:  "tcp",
					Url:      "localhost:6379",
					DBNumber: 0,
				}

				outputClusterRedisConfig{
					Url: []string{"localhost:6379", "localhost:6380"},
				}

	*/

	dump.Sources = append(dump.Sources, sourceConfig{
		File: sourceFileConfig{
			Filename:   "bk.resp",
			Batch:      dump.Common.Batch,
			ReaderSize: 4096 * 32,
		},
		Single: sourceSingleRedisConfig{
			Network:  "tcp",
			Url:      "localhost:6379",
			DBNumber: 0,
			TTL:      true,
			Match:    "*",
		},
		Cluster: sourceClusterRedisConfig{
			Url:   []string{"localhost:6379"},
			TTL:   true,
			Match: "*",
		},
	})

	dump.Outputs = append(dump.Outputs, outputConfig{
		File: outputFileConfig{
			Flag:      "trunc",
			Filename:  "test.resp",
			WriteSize: 4096 * 32,
		},
		Single: outputSingleRedisConfig{
			Network:  "tcp",
			Url:      "localhost:6379",
			DBNumber: 0,
		},
		Cluster: outputClusterRedisConfig{
			Url: []string{"localhost:6379", "localhost:6380"},
		},
	})

	buf, _ := json.Marshal(dump)
	fmt.Println(string(buf))
}
