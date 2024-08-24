package v2

import (
	"github.com/askolesov/image-vault/pkg/util"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

type FileWithSidecars struct {
	Path     string
	Sidecars []string
}

func LinkSidecars(
	sidecarExtensions []string,
	files []string,
) []FileWithSidecars {
	// 1. split files into primaries and sidecars

	sidecarExts := lo.Associate(sidecarExtensions, func(item string) (string, any) {
		return strings.ToLower(item), true
	})

	isSidecar := func(f string) bool {
		_, ok := sidecarExts[strings.ToLower(filepath.Ext(f))]
		return ok
	}

	var primaries, sidecars []string

	for _, f := range files {
		if isSidecar(f) {
			sidecars = append(sidecars, f)
		} else {
			primaries = append(primaries, f)
		}
	}

	// 2. add primaries with their sidecars to the result

	var result []FileWithSidecars

	sidecarsByPathWithoutExt := lo.GroupBy(sidecars, PathWithoutExtension)

	for _, f := range primaries {
		fs := FileWithSidecars{
			Path: f,
		}

		pathWithoutExt := util.GetPathWithoutExtension(f)
		if fSidecars, ok := sidecarsByPathWithoutExt[pathWithoutExt]; ok {
			fs.Sidecars = fSidecars
		}

		result = append(result, fs)
	}

	// 3. add sidecars without primaries to the result

	primariesByPathWithoutExt := lo.GroupBy(primaries, PathWithoutExtension)

	for _, f := range sidecars {
		pathWithoutExt := util.GetPathWithoutExtension(f)
		if _, ok := primariesByPathWithoutExt[pathWithoutExt]; !ok {
			result = append(result, FileWithSidecars{
				Path: f,
			})
		}
	}

	return result
}

func PathWithoutExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
