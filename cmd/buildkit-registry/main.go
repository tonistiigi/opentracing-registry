package main

import (
	"os"

	"github.com/moby/buildkit/client/llb"
	gobuild "github.com/tonistiigi/llb-gobuild"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	src := llb.Local("src")

	gb := gobuild.New(nil)

	registry, err := gb.BuildExe(gobuild.BuildOpt{
		Source:    src,
		MountPath: "/go/src/github.com/hinshun/opentracing-registry",
		Pkg:       "github.com/hinshun/opentracing-registry/cmd/registry",
		BuildTags: []string{},
	})
	if err != nil {
		return err
	}

	sc := llb.Scratch().
		With(copyAll(*registry, "/"))

	dt, err := sc.Marshal()
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)
	return nil
}

func copyAll(src llb.State, destPath string) llb.StateOption {
	return copyFrom(src, "/.", destPath)
}

// copyFrom has similar semantics as `COPY --from`
func copyFrom(src llb.State, srcPath, destPath string) llb.StateOption {
	return func(s llb.State) llb.State {
		return copy(src, srcPath, s, destPath)
	}
}

// copy copies files between 2 states using cp until there is no copyOp
func copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	cpImage := llb.Image("docker.io/library/alpine@sha256:1072e499f3f655a032e88542330cf75b02e7bdf673278f701d7ba61629ee3ebe")
	cp := cpImage.Run(llb.Shlexf("cp -a /src%s /dest%s", srcPath, destPath))
	cp.AddMount("/src", src, llb.Readonly)
	return cp.AddMount("/dest", dest)
}
