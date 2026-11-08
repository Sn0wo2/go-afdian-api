package afdian

import (
	"path"
	"runtime/debug"
	"sync"
)

var userAgent = sync.OnceValue(func() string {
	const modulePath = "github.com/Sn0wo2/go-afdian-api"

	name := path.Base(modulePath)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return name + "/dev"
	}

	v := info.Main.Version
	if info.Main.Path != modulePath {
		v = ""

		for _, dep := range info.Deps {
			if dep.Path != modulePath {
				continue
			}

			if dep.Replace != nil {
				v = dep.Replace.Version
			} else {
				v = dep.Version
			}

			break
		}
	}

	if v == "" || v == "(devel)" {
		return name + "/dev"
	}

	return name + "/" + v
})
