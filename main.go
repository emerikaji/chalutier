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

	useDigest :=
		flag.Bool("d", false, "Utiliser un digest à la place d'un id")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Command : chalutier (<docker container id> | -d <image digest>) [-r <docker registry>]")

		os.Exit(0)
	}

	digest := flag.Args()[0]

	if !*useDigest {
		digest =
			Catch(digests.ContainerDigest(flag.Args()[0]))
	}

	fmt.Print(Catch(digests.DigestVersion(*containerRegistry, digest)))
}

// ─────────────────────────────────────────────────────────────────────────────
