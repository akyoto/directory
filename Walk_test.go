package directory_test

import (
	"os"
	"testing"

	"github.com/akyoto/directory"
)

func BenchmarkWalk(b *testing.B) {
	callBack := func(name string) {}

	for i := 0; i < b.N; i++ {
		directory.Walk(".", callBack)
	}
}

func BenchmarkReaddirnames(b *testing.B) {
	callBack := func(name string) {}

	for i := 0; i < b.N; i++ {
		file, err := os.Open(".")

		if err != nil {
			b.Fatal(err)
		}

		files, err := file.Readdirnames(0)

		if err != nil {
			file.Close()
			b.Fatal(err)
		}

		for _, name := range files {
			callBack(name)
		}

		file.Close()
	}
}
