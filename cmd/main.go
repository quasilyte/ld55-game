package main

import (
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var g gameRunner
	g.ctx = game.NewContext()

	assets.RegisterResources(g.ctx.Loader)
	g.ctx.UIResources = eui.LoadResources(g.ctx.Loader)

	ebiten.SetWindowTitle("Astro Heart")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(int(g.ctx.WindowSize.X), int(g.ctx.WindowSize.Y))

	game.ChangeScene(g.ctx, scenes.NewMainMenuController(g.ctx))

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}

type gameRunner struct {
	ctx *game.Context
}

func (g *gameRunner) Update() error {
	const delta = 1.0 / 60.0
	g.ctx.AudioSystem.Update()
	g.ctx.CurrentScene().UpdateWithDelta(delta)

	return nil
}

func (g *gameRunner) Draw(screen *ebiten.Image) {
	g.ctx.CurrentScene().Draw(screen)
}

func (g *gameRunner) Layout(_, _ int) (int, int) {
	panic("unreachable")
}

func (g *gameRunner) LayoutF(outWidth, outHeight float64) (float64, float64) {
	return g.ctx.WindowSize.X, g.ctx.WindowSize.Y
}
