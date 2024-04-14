package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

func registerShaderResources(loader *resource.Loader) {
	resources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderCRT: {Path: "shader/crt.go"},
	}

	for id, res := range resources {
		loader.ShaderRegistry.Set(id, res)
		loader.LoadShader(id)
	}
}

const (
	ShaderNone resource.ShaderID = iota

	ShaderCRT
)
