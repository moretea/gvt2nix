package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type ManifestDependency struct {
	ImportPath string `json:"importpath"`
	Repository string `json:"repository"`
	Vcs        string `json:"vcs"`
	Revision   string `json:"revision"`
	Branch     string `json:"branch"`
	Path       string `json:"path"`
	NoTests    bool   `json:"notests"`
}

type Manifest struct {
	Version      int                  `json:"version"`
	Dependencies []ManifestDependency `json:"dependencies"`
}

type NixDependency struct {
	Dependency ManifestDependency `json:"dependency"`
	Sha256     string             `json:"sha256"`
}

func readManifest(path string) Manifest {
	rawJSON, err := ioutil.ReadFile(path)

	if err != nil {
		panic(fmt.Sprintf("Could not read %v!", err))
	}

	dat := Manifest{}
	err = json.Unmarshal(rawJSON, &dat)

	if err != nil {
		panic(err)
	}

	return dat
}

func processDependency(dep ManifestDependency) NixDependency {
	if dep.Vcs != "git" {
		panic("Only git dependencies are supported for now")
	}

	fmt.Printf("- Processing %s\n", dep.ImportPath)

	// Prefetch git
	prefetchCommand := exec.Command("nix-prefetch-git", "--rev", dep.Revision, "--no-deepClone", dep.Repository)
	prefetchJSON, err := prefetchCommand.Output()

	if err != nil {
		panic(fmt.Sprintf("Could not prefetch git repository! %v", err))
	}

	var nixPrefetchOutput map[string]interface{}
	err = json.Unmarshal(prefetchJSON, &nixPrefetchOutput)
	if err != nil {
		panic(fmt.Sprintf("parsing nix-prefetch-git output failed! %v", err))
	}

	sha256 := nixPrefetchOutput["sha256"].(string)

	if len(sha256) != 52 {
		panic("Wrong SHA256 in nix-prefetch-git output")
	}

	nixDependency := NixDependency{
		Dependency: dep,
		Sha256:     sha256,
	}

	return nixDependency
}

func main() {

	var manifestPath string
	var outputPath string

	flag.StringVar(&manifestPath, "manifest-path", "vendor/manifest", "path to the gvt manifest file")
	flag.StringVar(&outputPath, "output-path", "gvt2nix.json", "path to the gvt2nix dependency information file")
	flag.Parse()

	manifest := readManifest(manifestPath)

	nixDependencies := []NixDependency{}

	for _, dep := range manifest.Dependencies {
		nixDependencies = append(nixDependencies, processDependency(dep))
	}

	rawJSON, err := json.Marshal(nixDependencies)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing output file %s\n", outputPath)
	ioutil.WriteFile(outputPath, rawJSON, 0644)
}
