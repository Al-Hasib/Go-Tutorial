package main

import (
	"bytes"
	"log"
	"log/slog"
	"os"
)

func main() {

	// ---------- log: THE BASIC PACKAGE ----------

	//the standard "log" package prints a timestamp + message by default -
	//fine for a small CLI tool, limited once you need structure (fields
	//you can filter/search on, not just free text)
	log.Println("server starting up")
	log.Printf("listening on port %d", 8080)

	//a custom *log.Logger lets you set your own prefix/flags, or write
	//somewhere other than stderr (the default)
	var buf bytes.Buffer
	customLogger := log.New(&buf, "[worker] ", log.Ltime|log.Lshortfile)
	customLogger.Println("processed a job")
	os.Stdout.WriteString(buf.String())

	//log.Fatal(f/ln) logs, then calls os.Exit(1) IMMEDIATELY - deferred
	//functions do NOT run. that makes it a poor fit for anywhere with
	//cleanup to do (closing files, flushing a buffer, releasing a lock);
	//it's better suited to an early, nothing-set-up-yet startup failure.
	//not called here on purpose, since it would end this lesson early:
	//   log.Fatal("could not open config file") // -> os.Exit(1), no defers run

	// ---------- log/slog: STRUCTURED LOGGING ----------

	//slog (standard library since Go 1.21) attaches structured key-value
	//pairs to each entry instead of just a free-text message - far easier
	//to filter/search/aggregate in a real logging system than parsing text
	slog.Info("user logged in", "userID", 42, "method", "password")
	slog.Warn("rate limit approaching", "remaining", 5, "resetsIn", "30s")
	slog.Error("payment failed", "orderID", "ord_123", "reason", "card declined")

	//slog.With() returns a logger that always includes some fields - handy
	//for a "logger scoped to this request/job" without repeating fields on
	//every call
	requestLogger := slog.With("requestID", "req_abc123")
	requestLogger.Info("handling request")
	requestLogger.Info("request completed", "status", 200)

	//by default slog prints human-readable text - swapping the Handler
	//changes the OUTPUT FORMAT without touching any of the call sites above
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	jsonLogger.Info("this line is JSON instead of text", "userID", 42)

	//slog.SetDefault installs a logger for the whole program to use via the
	//package-level slog.Info/Warn/Error functions, exactly like the calls
	//above - useful for switching a whole app's log FORMAT in one place
	slog.SetDefault(jsonLogger)
	slog.Info("now every slog.Info call anywhere in the program is JSON")

	// ---------- LOG LEVELS ----------

	//a Handler can be given a minimum level - anything below it is dropped
	//before it's even formatted, which matters for performance in a hot path
	quietLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn, // Info and Debug calls below this are silently skipped
	}))
	quietLogger.Info("you will NOT see this - below the configured level")
	quietLogger.Warn("you WILL see this - at or above the configured level")
}
