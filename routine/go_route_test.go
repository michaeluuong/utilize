package routine

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)

}

func setup() {

}

func teardown() {

}

// go test -benchtime=1s -bench . -cpuprofile cpu.prof
// go tool pprof cpu.prof
func BenchmarkReflections(b *testing.B) {

}

func TestAnyType_main(t *testing.T) {
	//--------------------------------------------------------------------------------------------
	//--------------------------------------------------------------------------------------------
	// Cleanup

}
