package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/emerikaji/chalutier/digests"
)

// ─── Generic Error Throw ─────────────────────────────────────────────────────

func Catch[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return value
}

// ─── Main ────────────────────────────────────────────────────────────────────

func main() {
	containerRegistry :=
		flag.String("r", "https://hub.docker.com", "Docker registry")

	flag.Parse()

	if len(flag.Args()) > 2 {
		fmt.Println("Command : chalutier <docker container id> [-r <docker registry>]")

		os.Exit(0)
	}

	containerId :=
		flag.Args()[0]

	container :=
		Catch(digests.Inspect(containerId))

	Catch("", container.GetVersion(*containerRegistry))

	fmt.Print(container.Version)
}

// ─────────────────────────────────────────────────────────────────────────────
