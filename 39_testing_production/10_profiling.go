package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* on http.DefaultServeMux, purely as an import side effect
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

func main() {

	profileDir := "profiles" // relative to this lesson folder

	// ---------- CPU PROFILING ----------

	//StartCPUProfile samples the call stack periodically while it's
	//running, writing samples to the given file - StopCPUProfile flushes
	//and finishes the file
	cpuFile, err := os.Create(filepath.Join(profileDir, "cpu.prof"))
	must(err)
	defer cpuFile.Close()

	must(pprof.StartCPUProfile(cpuFile))
	fmt.Println("profiling some CPU-heavy work...")
	//the CPU profiler samples the call stack ~100 times/second - fibonacci(30)
	//finishes in under 10ms, too fast to collect more than a sample or two;
	//38 runs long enough (a few hundred ms) to give the profiler something
	//real to find, which matters for this lesson actually showing anything
	result := slowFibonacci(38)
	pprof.StopCPUProfile()
	fmt.Println("fibonacci(38) =", result)
	fmt.Println("wrote", cpuFile.Name())

	// ---------- MEMORY (HEAP) PROFILING ----------

	allocateLotsOfMemory()
	runtime.GC() // a heap profile taken right after a GC reflects live objects, not garbage awaiting collection

	memFile, err := os.Create(filepath.Join(profileDir, "mem.prof"))
	must(err)
	defer memFile.Close()
	must(pprof.WriteHeapProfile(memFile))
	fmt.Println("wrote", memFile.Name())

	fmt.Println("\nanalyze either file with:")
	fmt.Println("  go tool pprof", filepath.Join(profileDir, "cpu.prof"))
	fmt.Println("  go tool pprof", filepath.Join(profileDir, "mem.prof"))
	fmt.Println("then at its (pprof) prompt: `top` for the biggest consumers, or `web` for a call graph (needs graphviz)")

	// ---------- net/http/pprof: THE PRODUCTION PATTERN ----------

	//blank-importing net/http/pprof above registers /debug/pprof/* handlers
	//as a side effect of the import alone - no explicit setup code needed.
	//the usual real-world setup serves that mux (often on a separate,
	//non-public port) so a LIVE, already-running service can be profiled
	//on demand, instead of only ever profiling one fixed snippet like above.
	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/debug/pprof/")
	must(err)
	fmt.Println("\nGET /debug/pprof/ on a live server -> status:", resp.StatusCode)
	resp.Body.Close()
	fmt.Println("(in production you'd point at this over the network: go tool pprof http://host:port/debug/pprof/profile)")
}

func slowFibonacci(n int) int {
	if n < 2 {
		return n
	}
	return slowFibonacci(n-1) + slowFibonacci(n-2)
}

func allocateLotsOfMemory() {
	var chunks [][]byte
	for i := 0; i < 1000; i++ {
		chunks = append(chunks, make([]byte, 10_000))
	}
	_ = chunks
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
