package main

import (
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
)

func main() {
	out := registry().Run(llb.Shlex("ls -l /bin")) // debug output

	dt, err := out.Marshal()
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)
}

func goBuildBase() llb.State {
	goAlpine := llb.Image("docker.io/library/golang:1.8-alpine")
	return goAlpine.
		AddEnv("PATH", "/usr/local/go/bin:"+system.DefaultPathEnv).
		AddEnv("GOPATH", "/go").
		AddEnv("GOOS", "linux").
		AddEnv("GOARCH", "amd64").
		Run(llb.Shlex("apk add --no-cache git make")).Root()
}

func registry() llb.State {
	src := goBuildBase().
		Run(llb.Shlex("git clone https://github.com/hinshun/opentracing-registry.git /go/src/github.com/hinshun/opentracing-registry")).
		Dir("/go/src/github.com/hinshun/opentracing-registry")

	registry := src.
		Run(llb.Shlex("go build -o /bin/registry ./cmd/registry"))

	r := llb.Image("docker.io/library/golang:1.8-alpine")
	r = copy(registry.Root(), "/bin/registry", r, "/bin/")
	return r
}

func copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	cpImage := llb.Image("docker.io/library/alpine:latest")
	cp := cpImage.Run(llb.Shlexf("cp -a /src%s /dest%s", srcPath, destPath))
	cp.AddMount("/src", src)
	return cp.AddMount("/dest", dest)
}
