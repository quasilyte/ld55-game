package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/gameui/eui"
)

type Context struct {
	WindowSize gmath.Vec

	Loader *resource.Loader

	UIResources *eui.Resources

	GraphicsCache *graphics.Cache

	Rand gmath.Rand

	scene gscene.GameRunner
}

func ChangeScene[C any](ctx *Context, c gscene.Controller[C]) {
	s := gscene.NewRootScene[C](c)
	ctx.scene = s
}

func NewContext() *Context {
	sampleRate := 44100
	audioContext := audio.NewContext(sampleRate)
	l := resource.NewLoader(audioContext)
	l.OpenAssetFunc = assets.OpenAssetFunc
	ctx := &Context{
		WindowSize: gmath.Vec{
			X: 1920 / 2,
			Y: 1080 / 2,
		},
		Loader:        l,
		GraphicsCache: graphics.NewCache(),
	}
	ctx.Rand.SetSeed(time.Now().Unix())
	return ctx
}

func (ctx *Context) CurrentScene() gscene.GameRunner {
	return ctx.scene
}

func (ctx *Context) NewRect(w, h float64) *graphics.Rect {
	r := graphics.NewRect(ctx.GraphicsCache, w, h)
	return r
}

func (ctx *Context) NewSprite(id resource.ImageID) *graphics.Sprite {
	s := graphics.NewSprite(ctx.GraphicsCache)
	img := ctx.Loader.LoadImage(id)
	s.SetImage(img.Data)
	if img.DefaultFrameWidth != 0 {
		_, fh := s.GetFrameSize()
		s.SetFrameSize(int(img.DefaultFrameWidth), fh)
	}
	return s
}
