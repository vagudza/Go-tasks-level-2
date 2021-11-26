package main

import "testing"

func BenchmarkWget(b *testing.B) {
	conf := Config{
		url:          "https://www.iana.org/",
		downloadPath: "downloads",
		maxDepth:     1,
		delay:        0,
	}
	wget := NewWget(&conf)

	for i := 0; i < b.N; i++ {
		wget.wget(conf.url, conf.downloadPath, conf.maxDepth)
	}
}

// go test -bench=.
// GRAPHVIZ работает только из Powershell (не из vscode)
// go tool pprof -http=:8080 cpu.out
