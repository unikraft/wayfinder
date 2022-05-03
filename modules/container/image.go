package container

// SPDX-License-Identifier: Apache-2.0
//
// Copyright (C) 2015-2017 Thomas LE ROUX <thomas@leroux.io>
//               2020      Alexander Jung <a.jung@lancs.ac.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Parse parses s and returns a syntactically valid Reference.
// If an error was encountered it is returned, along with a nil Reference.
// NOTE: Parse will not handle short digests.

import (
	"encoding/json"
	// "errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/docker/docker/daemon/graphdriver/copy"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/moby/moby/pkg/archive"
	"github.com/tidwall/gjson"
)

// const (
// 	// DefaultRuntime is the runtime to use when not specified.
// 	DefaultRuntime = "runc"
// 	// NameTotalLengthMax is the maximum total number of characters in a
// 	// repository name.
// 	NameTotalLengthMax = 255
// 	// DefaultTag defines the default tag used when performing images related
// 	// actions and no tag or digest is specified.
// 	DefaultTag = "latest"
// 	// DefaultHostname is the default built-in hostname
// 	DefaultHostname = "docker.io"
// 	// LegacyDefaultHostname is automatically converted to DefaultHostname
// 	LegacyDefaultHostname = "index.docker.io"
// 	// DefaultRepoPrefix is the prefix used for default repositories in default
// 	// host.
// 	DefaultRepoPrefix = "library/"
// )

// var (
// 	// ErrReferenceInvalidFormat represents an error while trying to parse a
// 	// string as a reference.
// 	ErrReferenceInvalidFormat = errors.New("invalid reference format")
// 	// ErrTagInvalidFormat represents an error while trying to parse a string as a
// 	// tag.
// 	ErrTagInvalidFormat = errors.New("invalid tag format")
// 	// ErrDigestInvalidFormat represents an error while trying to parse a string
// 	// as a tag.
// 	ErrDigestInvalidFormat = errors.New("invalid digest format")
// 	// ErrNameEmpty is returned for empty, invalid repository names.
// 	ErrNameEmpty = errors.New("repository name must have at least one component")
// 	// ErrNameTooLong is returned when a repository name is longer than
// 	// NameTotalLengthMax.
// 	ErrNameTooLong = fmt.Errorf("repository name must not be more than %v characters", NameTotalLengthMax)
// )

// // Image is an object with a full name
// type Image struct {
// 	// Runtime is the normalized name of the runtime service, e.g. "docker"
// 	Runtime string
// 	// Name is the normalized repository name, like "ubuntu".
// 	Name string
// 	// String is the full reference, like "ubuntu@sha256:abcdef..."
// 	String string
// 	// FullName is the full repository name with hostname, like "docker.io/library/ubuntu"
// 	FullName string
// 	// Hostname is the hostname for the reference, like "docker.io"
// 	Hostname string
// 	// RemoteName is the the repository component of the full name, like "library/ubuntu"
// 	RemoteName string
// 	// Tag is the tag of the image, e.g. "latest"
// 	Tag string
// }

// Image type to save the provider
type Image struct {
	P *Provider
}

// PullImage downloads an image
func (i *Image) PullImage(image, cacheDir, savedDir string) (v1.Image, string, error) {
	var options []crane.Option

	// options = append(options, crane.Insecure)

	// Use current built OS and architecture
	options = append(options, crane.WithPlatform(&v1.Platform{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}))

	// Grab the remote manifest
	manifest, err := crane.Manifest(image, options...)
	if err != nil {
		return nil, "", fmt.Errorf("failed fetching manifest for %s: %v", image, err)
	}

	if !gjson.Valid(string(manifest)) {
		return nil, "", fmt.Errorf("cannot parse manifest: %s", string(manifest))
	}

	value := gjson.Get(string(manifest), "config.digest").Value().(string)
	if value == "" {
		return nil, "", fmt.Errorf("malformed manifest: %s", string(manifest))
	}

	digest := strings.Split(value, ":")[1]
	tarball := fmt.Sprintf("%s/%s.tar.gz", cacheDir, digest)

	// Download the tarball of the image if not available in the cache
	if _, err := os.Stat(tarball); os.IsNotExist(err) {
		// Create the cacheDir if it does not already exist
		if cacheDir != "" {
			if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
				os.MkdirAll(cacheDir, os.ModePerm)
			}
		}

		// Pull the image
		img, err := crane.Pull(image, options...)
		if err != nil {
			return nil, "", fmt.Errorf("could not pull image: %s", err)
		}

		f, err := os.Create(tarball)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open %s: %v", tarball, err)
		}

		defer f.Close()

		// Save the image permanently in the registry
		serverAddress := os.Getenv("REGISTRY_ADDR")
		splitPath := strings.Split(image, "/")
		name := splitPath[len(splitPath)-1]

		err = crane.Push(img, serverAddress+"/"+name, options...)
		if err != nil {
			return nil, "", fmt.Errorf("failed to push image: %s", err)
		}

		err = crane.Save(img, image, tarball)
		if err != nil {
			return nil, "", fmt.Errorf("could not save image: %s", err)
		}
	}

	var img v1.Image = nil
	imageDir := fmt.Sprintf("%s/%s", savedDir, digest)
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		img, err = crane.Load(tarball)
		if err != nil {
			return nil, "", fmt.Errorf("could not load image: %s", err)
		}
	}

	return img, digest, nil
}

// Return the configuration for an image
func (i *Image) ImageConfig(image string) (*v1.ConfigFile, error) {
	var options []crane.Option

	// Use current built OS and architecture
	options = append(options, crane.WithPlatform(&v1.Platform{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}))

	// Grab the remote manifest
	byt, err := crane.Config(image, options...)
	if err != nil {
		return nil, fmt.Errorf("failed fetching manifest for %s: %v", image, err)
	}

	var cfg v1.ConfigFile
	if err := json.Unmarshal(byt, &cfg); err != nil {
		return nil, fmt.Errorf("could not serialize v1.Config: %s", err)
	}

	return &cfg, nil
}

// UnpackImage takes a container image and writes its filesystem to outDir
func (i *Image) UnpackImage(image v1.Image, manifestHex string, cacheDir, savedDir, outDir string, allowOverride bool) error {
	imageDir := fmt.Sprintf("%s/%s", savedDir, manifestHex)

	// Check if the tarball exists in the cache
	tarball := fmt.Sprintf("%s/%s.tar.gz", cacheDir, manifestHex)
	if _, err := os.Stat(tarball); os.IsNotExist(err) {
		return fmt.Errorf("image does not exist: %s", tarball)
	}

	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		if image == nil {
			return fmt.Errorf("invalid image")
		}

		// Get the layers of this image
		layers, err := image.Layers()
		if err != nil {
			return fmt.Errorf("could not read layers: %s", err)
		}

		for _, layer := range layers {
			uncompressed, err := layer.Uncompressed()
			if err != nil {
				return fmt.Errorf("could not read layer: %s", err)
			}

			// Untar cached tarball
			err = archive.Untar(uncompressed, imageDir, &archive.TarOptions{
				NoLchown: true,
			})
			if err != nil {
				return fmt.Errorf("extracting tar from %s failed: %s", tarball, err)
			}
		}
	}

	// Copies the saved content to the output directory
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.MkdirAll(outDir, os.ModePerm)
	}
	err := copy.DirCopy(imageDir, outDir, copy.Content, false)
	if err != nil {
		return fmt.Errorf("could not copy: %s", err)
	}

	return nil
}
