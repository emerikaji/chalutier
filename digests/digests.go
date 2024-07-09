package digests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// ─── Commands ────────────────────────────────────────────────────────────────

const (
	containerInspect = "docker container inspect --format '{{.Image}}' %s"
	imageInspect     = "docker image inspect --format '{{ join .RepoDigests \",\"}}' %s"
)

// ─── Structs ─────────────────────────────────────────────────────────────────

type DockerContainer struct {
	id        string
	name      string
	namespace string
	digest    string
	Version   string
}

type DockerTags struct {
	NextUrl *string `json:"next"`
	Results []struct {
		Id      int    `json:"id"`
		Version string `json:"name"`
		Digest  string `json:"digest"`
	} `json:"results"`
}

// ─── Compare Digest With Tags ────────────────────────────────────────────────

func (d *DockerContainer) GetVersion(registry string) (err error) {
	if d.name == "" || d.digest == "" {
		return
	}

	next :=
		fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/tags?page_size=100", registry, d.namespace, d.name)

	tags :=
		&DockerTags{
			NextUrl: &next,
		}

	for tags.NextUrl != nil {
		var resp *http.Response

		if resp, err = http.Get(*tags.NextUrl); err != nil {
			return
		}

		if err = json.NewDecoder(resp.Body).Decode(tags); err != nil {
			return
		}

		resp.Body.Close()

		// ─── Compare Tag And Image Digests ───────────────────────────

		var regex *regexp.Regexp

		if regex, err = regexp.Compile(".?[0-9].*"); err != nil {
			return
		}

		for _, tag := range tags.Results {
			if d.digest == tag.Digest {
				if regex.Match([]byte(tag.Version)) {
					d.Version = tag.Version
					return
				}
			}
		}
	}

	return
}

// ─── Inspect ─────────────────────────────────────────────────────────────────

func Inspect(id string) (container *DockerContainer, err error) {
	imageIdBytes, err :=
		exec.Command("bash", "-c", fmt.Sprintf(containerInspect, id)).Output()

	if err != nil {
		return
	}

	imageDigestsBytes, err :=
		exec.Command("bash", "-c", fmt.Sprintf(imageInspect, string(imageIdBytes))).Output()

	if err != nil {
		return
	}

	imageDigest :=
		strings.Split(strings.ReplaceAll(string(imageDigestsBytes), "\n", ""), ",")[0]

	container = &DockerContainer{
		id:      id,
		Version: "null",
	}

	container.namespace, container.name, container.digest =
		parseDigest(imageDigest)

	return
}

// ─── Parse Digest String ─────────────────────────────────────────────────────

func parseDigest(fullDigest string) (namespace string, name string, digest string) {
	fullName, digest, ok :=
		strings.Cut(strings.Split(strings.ReplaceAll(fullDigest, "\n", ""), ",")[0], "@")

	if !ok {
		digest = fullDigest
	}

	namespace, name, ok =
		strings.Cut(fullName, "/")

	if !ok {
		namespace = "library"
		name = fullName
	}

	return
}

// ─────────────────────────────────────────────────────────────────────────────
