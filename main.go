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
var checksum = flag.Bool("checksum", true, "Generate sha256 checksum")
var appName = flag.String("appName", "poco App name", "Application name")
var appCopyright = flag.String("appCopyright", "noone", "Application copyright")
var appAuthor = flag.String("appAuthor", "noone", "Application author")
var appVersion = flag.String("appVersion", "0.1", "Application version")
var appMounts = flag.String("appMounts", "/etc/resolv.conf", "Application mounts")
var appAttrs = flag.String("appAttrs", "ipc,uts,user,ns,pid", "Application attrs")
var appStore = flag.String("appStore", ".store", "Application store")
var image = flag.String("image", "", "Application image")
var compression = flag.String("compression", "xz", "Application compression")

var version = flag.String("version", "v0.2.1", "poco version")
var arch = flag.String("arch", "x86_64", "poco architecture")

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

	RunSH("dependencies", fmt.Sprintf("curl -L https://github.com/mudler/poco/releases/download/%s/poco-%s-Linux-%s.tar.gz --output poco.tar.gz", *version, *version, *arch))
	RunSH("dependencies", "tar xvf poco.tar.gz")
	RunSH("dependencies", "mv poco /usr/bin/poco")
	RunSH("dependencies", "poco --version")

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
	
	if *checksum {
		checkErr(RunSH("sha256sum", fmt.Sprintf("sha256sum %s > %s.sha256",*output,*output)))
	}
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
