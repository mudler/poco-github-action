package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var entrypoint = flag.String("appEntrypoint", "/bin/sh", "Application container entrypoint")
var output = flag.String("output", "output", "Output file path")
var appDescription = flag.String("appDescription", "poco Description", "Application description")
var localDaemon = flag.Bool("local", true, "Use local daemon to find image")
var appName = flag.String("appName", "poco App name", "Application name")
var appCopyright = flag.String("appCopyright", "", "Application copyright")
var appAuthor = flag.String("appAuthor", "", "Application author")
var appVersion = flag.String("appVersion", "", "Application version")
var appMounts = flag.String("appMounts", "", "Application mounts")
var appAttrs = flag.String("appAttrs", "", "Application attrs")
var appStore = flag.String("appStore", ".store", "Application store")
var image = flag.String("image", "", "Application image")
var compression = flag.String("compression", "", "Application compression")

var version = flag.String("version", "v0.2.1", "poco version")

func RunSH(stepName, bashFragment string, envs ...string) error {
	cmd := exec.Command("sh", "-s")
	cmd.Stdin = strings.NewReader(bashWrap(bashFragment))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), envs...)
	log.Printf("Running: %v (%v)", stepName, bashFragment)

	return cmd.Run()
}

func bashWrap(cmd string) string {
	return `
set -o errexit
set -o nounset
` + cmd + `
`
}

func main() {
	flag.Parse()

	RunSH("dependencies", fmt.Sprintf("curl -L https://github.com/mudler/poco/releases/download/%s/poco-%s-Linux-x86_64.tar.gz --output poco.tar.gz", *version, *version))
	RunSH("dependencies", "tar xvf poco.tar.gz")
	RunSH("dependencies", "sudo mv poco /usr/bin/poco")
	l := "--local"
	if !*localDaemon {
		l = ""
	}
	checkErr(RunSH("build", fmt.Sprintf("CGO_ENABLED=0 poco bundle --command-prefix '' --entrypoint %s --output %s %s --image %s --compression %s", *entrypoint, *output, l, *image, *compression),
		toKey("DESCRIPTION", *appDescription),
		toKey("COPYRIGHT", *appCopyright),
		toKey("AUTHOR", *appAuthor),
		toKey("NAME", *appName),
		toKey("VERSION", *appVersion),
		toKey("MOUNTS", *appMounts),
		toKey("ATTRS", *appAttrs),
		toKey("STORE", *appStore),
	))
}

func toKey(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
