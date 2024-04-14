package assets

import (
	"embed"
	"io"

	resource "github.com/quasilyte/ebitengine-resource"
)

//go:embed all:_data
var data embed.FS

func OpenAssetFunc(p string) io.ReadCloser {
	fullPath := "_data/" + p
	r, err := data.Open(fullPath)
	if err != nil {
		panic(err)
	}
	return r
}

func RegisterResources(loader *resource.Loader) {
	registerImageResources(loader)
	registerAudioResources(loader)
	registerRawResources(loader)
	registerShaderResources(loader)
}
