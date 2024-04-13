package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/gameui/eui"
)

func initUI(scene *gscene.SimpleRootScene, root *widget.Container) {
	uiObject := eui.NewSceneObject(root)
	scene.AddObject(uiObject)
	scene.AddGraphics(uiObject)
}
