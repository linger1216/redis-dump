package core

import (
	"fmt"
	"testing"
)

func TestStandAloneRedisSource(t *testing.T) {

	var r source
	var w output

	r = newSourceFile(&sourceFileConfig{
		Filename:   "local.resp",
		Batch:      3,
		ReaderSize: 4096,
	})

	//r := &sourceSingleRedis{
	//	c: redis.NewClient(&redis.Options{
	//		Addr: "localhost:6379",
	//	}),
	//	TTL:   true,
	//	Match: "*",
	//	Count: 100,
	//}

	//w := NewOutputSingleRedis(&OutputSingleRedisConfig{
	//	Url: "localhost:7379",
	//})

	w = NewOutputFile(&OutputFileConfig{
		Flag:      "trunc",
		Filename:  "local1.resp",
		WriteSize: 4096,
	})

	//w := NewOutputSingleRedis(&OutputSingleRedisConfig{
	//	Url:      "localhost:7379",
	//})

	//w := NewOutputClusterRedis(&OutputClusterRedisConfig{
	//	Url:      []string{"localhost:7379"},
	//})

	count := uint64(0)
	for r.has() {
		commands, err := r.next()
		if err != nil {
			t.Fatal(err)
		}
		err = w.save(commands)
		if err != nil {
			t.Fatal(err)
		}
		count += uint64(len(commands))
		fmt.Println(count)
	}

	r.close()
	w.close()
}

func TestRedisRestore(t *testing.T) {
	var r source
	var w output
	r = newSourceFile(&sourceFileConfig{
		Filename:   "../bk.resp",
		Batch:      3,
		ReaderSize: 4096,
	})

	//r := &sourceSingleRedis{
	//	c: redis.NewClient(&redis.Options{
	//		Addr: "localhost:6379",
	//	}),
	//	TTL:   true,
	//	Match: "*",
	//	Count: 100,
	//}

	w = NewOutputSingleRedis(&OutputSingleRedisConfig{
		Url: "localhost:6379",
	})
	//
	//w = NewOutputFile(&OutputFileConfig{
	//	Flag:      "trunc",
	//	Filename:  "local1.resp",
	//	WriteSize: 4096,
	//})

	//w := NewOutputSingleRedis(&OutputSingleRedisConfig{
	//	Url:      "localhost:7379",
	//})

	//w := NewOutputClusterRedis(&OutputClusterRedisConfig{
	//	Url:      []string{"localhost:7379"},
	//})

	count := uint64(0)
	for r.has() {
		commands, err := r.next()
		if err != nil {
			t.Fatal(err)
		}
		err = w.save(commands)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println(count)
	r.close()
	w.close()
}
