package utils

import (
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

type UserAgent struct{}

var userAgent = sync.OnceValue(func() string {
	moduleName := func(p string) string {
		segs := strings.Split(p, "/")

		for i, seg := range segs {
			if len(seg) > 1 && seg[0] == 'v' {
				if n, err := strconv.Atoi(seg[1:]); err == nil && n >= 2 {
					segs = segs[:i]

					break
				}
			}
		}

		return segs[len(segs)-1]
	}

	hasPathPrefix := func(s, prefix string) bool {
		return s == prefix || strings.HasPrefix(s, prefix+"/")
	}

	pkgPath := reflect.TypeOf(UserAgent{}).PkgPath()

	info, ok := debug.ReadBuildInfo()
	if ok {
		mod := &info.Main
		if !hasPathPrefix(pkgPath, mod.Path) {
			mod = nil

			for _, dep := range info.Deps {
				if hasPathPrefix(pkgPath, dep.Path) {
					mod = dep
					break
				}
			}
		}

		if mod != nil {
			v := mod.Version
			if mod.Replace != nil {
				v = mod.Replace.Version
			}

			if v != "" && v != "(devel)" {
				return moduleName(mod.Path) + "/" + v
			}
		}
	}

	return moduleName(pkgPath) + "/dev"
})

func (UserAgent) String() string {
	return userAgent()
}
