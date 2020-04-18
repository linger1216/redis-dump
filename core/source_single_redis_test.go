package core

//
//func TestStandAloneRedisSource(t *testing.T) {
//
//	var r source
//	var w output
//
//	r = newSourceFile(&sourceFileConfig{
//		Filename:   "local.resp",
//		Batch:      3,
//		ReaderSize: 4096,
//	})
//
//	//r := &sourceSingleRedis{
//	//	c: redis.NewClient(&redis.Options{
//	//		Addr: "localhost:6379",
//	//	}),
//	//	TTL:   true,
//	//	Match: "*",
//	//	Count: 100,
//	//}
//
//	//w := NewOutputSingleRedis(&outputSingleRedisConfig{
//	//	Url: "localhost:7379",
//	//})
//
//	w = NewOutputFile(&outputFileConfig{
//		Flag:      "trunc",
//		Filename:  "local1.resp",
//		WriteSize: 4096,
//	})
//
//	//w := NewOutputSingleRedis(&outputSingleRedisConfig{
//	//	Url:      "localhost:7379",
//	//})
//
//	//w := NewOutputClusterRedis(&outputClusterRedisConfig{
//	//	Url:      []string{"localhost:7379"},
//	//})
//
//	Count := uint64(0)
//	for r.has() {
//		commands, err := r.next()
//		if err != nil {
//			t.Fatal(err)
//		}
//		err = w.save(commands)
//		if err != nil {
//			t.Fatal(err)
//		}
//		Count += uint64(len(commands))
//		fmt.Println(Count)
//	}
//
//	r.close()
//	w.close()
//}
//
//func TestRedisRestore(t *testing.T) {
//	var r source
//	var w output
//	r = newSourceFile(&sourceFileConfig{
//		Filename:   "../bk.resp",
//		Batch:      3,
//		ReaderSize: 4096,
//	})
//
//	//r := &sourceSingleRedis{
//	//	c: redis.NewClient(&redis.Options{
//	//		Addr: "localhost:6379",
//	//	}),
//	//	TTL:   true,
//	//	Match: "*",
//	//	Count: 100,
//	//}
//
//	w = NewOutputSingleRedis(&outputSingleRedisConfig{
//		Url: "localhost:6379",
//	})
//	//
//	//w = NewOutputFile(&outputFileConfig{
//	//	Flag:      "trunc",
//	//	Filename:  "local1.resp",
//	//	WriteSize: 4096,
//	//})
//
//	//w := NewOutputSingleRedis(&outputSingleRedisConfig{
//	//	Url:      "localhost:7379",
//	//})
//
//	//w := NewOutputClusterRedis(&outputClusterRedisConfig{
//	//	Url:      []string{"localhost:7379"},
//	//})
//
//	Count := uint64(0)
//	for r.has() {
//		commands, err := r.next()
//		if err != nil {
//			t.Fatal(err)
//		}
//		err = w.save(commands)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	fmt.Println(Count)
//	r.close()
//	w.close()
//}
