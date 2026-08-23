package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config gathers everything the "app" needs to run, from wherever it
// actually comes from - callers of Load() don't need to know or care
// whether a value came from a flag, an env var, or a default.
type Config struct {
	Port    int
	Debug   bool
	DBHost  string
	Version string // hardcoded, not configurable - not everything needs to be
}

func (c Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	return nil
}

func main() {

	// ---------- WHY NOT JUST HARDCODE IT ----------

	//a hardcoded `const port = 8080` can't run two copies on different
	//ports, can't point at a staging vs. production database, and forces a
	//rebuild for a change that should just be a restart with new settings.

	// ---------- ENV VARS: THE USUAL BASELINE ----------

	//os.Getenv returns "" for a var that isn't set - there's no built-in
	//way to tell "unset" apart from "explicitly set to empty", so a helper
	//with an explicit default is the normal pattern
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	debug := getEnvBool("DEBUG", false)
	port := getEnvInt("PORT", 8080)

	// ---------- FLAGS: FOR VALUES SET AT LAUNCH, NOT VIA THE ENVIRONMENT ----------

	//flag.IntVar/BoolVar/StringVar let a value be set as `-port 9090` etc -
	//useful for local dev overrides ("just this once, use a different port")
	//without touching env vars at all. flag.Parse() reads os.Args.
	flagPort := flag.Int("port", port, "port to listen on (overrides PORT env var if set)")
	flagDebug := flag.Bool("debug", debug, "enable debug logging (overrides DEBUG env var if set)")
	flag.Parse()

	cfg := Config{
		Port:    *flagPort, // flag > env var > built-in default, in that precedence order
		Debug:   *flagDebug,
		DBHost:  dbHost,
		Version: "1.0.0",
	}

	// ---------- VALIDATE EARLY, FAIL LOUD ----------

	//catching a bad config at startup (before anything else runs) is far
	//better than discovering it 20 minutes later from a confusing runtime
	//error deep inside some unrelated code path
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid configuration:", err)
		//a real program would os.Exit(1) here - not done in this lesson so
		//it can keep running and show the rest below regardless
	}

	fmt.Printf("loaded config: %+v\n", cfg)
	fmt.Println("\n(run this file again as: go run 08_configuration.go -port=9090 -debug=true)")
	fmt.Println("(or with an env var: PORT=9090 DEBUG=true go run 08_configuration.go)")
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback // an unparseable value falls back rather than crashing config loading
	}
	return b
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
