package digests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/docker/docker/client"
)

// ─── Init ────────────────────────────────────────────────────────────────────

var dockerClient *client.Client

func init() {
	var err error

	dockerClient, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
}

// ─── Structs ─────────────────────────────────────────────────────────────────

type DockerTags struct {
	NextUrl *string `json:"next"`
	Results []struct {
		Version string `json:"name"`
		Digest  string `json:"digest"`
	} `json:"results"`
}

// ─── Container Version ───────────────────────────────────────────────────────

func ContainerDigest(id string) (fullDigest string, err error) {
	containerInfo, err := dockerClient.ContainerInspect(context.Background(), id)

	if err != nil {
		return
	}

	imageInfo, _, err := dockerClient.ImageInspectWithRaw(context.Background(), containerInfo.Image)

	if err != nil || len(imageInfo.RepoDigests) == 0 {
		return
	}

	return imageInfo.RepoDigests[0], nil
}

// ─── Compare Digest With Tags ────────────────────────────────────────────────

func DigestVersion(registry string, fullDigest string) (version string, err error) {
	version = "null"

	if fullDigest == "" {
		return
	}

	namespace, name, digest, err :=
		parseDigest(fullDigest)

	if err != nil {
		return
	}

	next :=
		fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/tags?page_size=100", registry, namespace, name)

	tags :=
		&DockerTags{
			NextUrl: &next,
		}

search:
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
			if digest == tag.Digest {
				if regex.Match([]byte(tag.Version)) {
					version = tag.Version
					break search
				}
			}
		}
	}

	return
}

// ─── Parse Digest String ─────────────────────────────────────────────────────

func parseDigest(fullDigest string) (namespace string, name string, digest string, err error) {
	fullName, digest, ok :=
		strings.Cut(fullDigest, "@")

	if !ok {
		return "", "", "", errors.New("incomplete digest")
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
