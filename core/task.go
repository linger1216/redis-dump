package core

import (
	"runtime"
)

func Exec(conf *DumpConfig) error {

	if conf.Common.Batch == 0 {
		conf.Common.Batch = 5000
	}

	//ch := make(chan [][]string, conf.Common.Batch)

	if conf.Common.Parallel == 0 {
		conf.Common.Parallel = runtime.NumCPU()
	}
	for i := range conf.Source.File {
		conf.Source.File[i].Batch = conf.Common.Batch
	}
	for i := range conf.Source.Single {
		conf.Source.Single[i].Count = int64(conf.Common.Batch)
	}
	for i := range conf.Source.Cluster {
		conf.Source.Cluster[i].Count = int64(conf.Common.Batch)
	}

	outputs := make([]output, 0, 3)

	if conf.Output.File != nil {
		if v := conf.Output.File.newOutput(); v != nil {
			outputs = append(outputs, v)
		}
	}

	if conf.Output.Single != nil {
		if v := conf.Output.Single.newOutput(); v != nil {
			outputs = append(outputs, v)
		}
	}

	if conf.Output.Cluster != nil {
		if v := conf.Output.Cluster.newOutput(); v != nil {
			outputs = append(outputs, v)
		}
	}

	read := func(s source, outputs []output) error {
		for s.has() {
			commands, err := s.next()
			if err != nil {
				return err
			}
			for _, w := range outputs {
				err = w.save(commands)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, v := range conf.Source.File {
		s := v.newSource()
		err := read(s, outputs)
		if err != nil {
			return err
		}
	}

	for _, v := range conf.Source.Single {
		s := v.newSource()
		err := read(s, outputs)
		if err != nil {
			return err
		}
	}

	for _, v := range conf.Source.Cluster {
		s := v.newSource()
		err := read(s, outputs)
		if err != nil {
			return err
		}
	}

	return nil
}
