LD55

Theme: summoning

+ Main Menu UI - 30
- Gameplay (out of combat) - 8:00
-+ Gameplay (in combat) - 3:00
- Sprites - 1:30
- SFX - 15
- Music - 30
- Testing - 1:00
- Final testing - 1:00
- Deploy - 1:00 (wasm)

TODO (top priority):

* software
  * give notification when action is failed due to an error
* journal

hardware:

* vessel design
* weapon 1
* weapon 2
* artifact (various bonuses)

TODO:

* weapon with slow condition
* homing rockets
* "enemy in arc" condition
* instruction progression (don't unlock all at level 1)
* menu bg
* weapon/vessel balancing
* vessel explosion effect (animation)
* cpu costs for different actions (snap/normal/snipe)
* sound level balancing (sfx)

TODO (polishing):

* optimize sfx (if needed)
* add input system
* add graphic layers
* add hotkeys support (at least "go back" at esc)
* button press sfx

Game loop:

* new game
  * pre-battle setup
  * battle
  * battle results

pre-battle:
- equipment setup
- drone program

battle:
- space battle with programmed vessels

# TODO engine

* add Animation to graphics pkg
* need audio sys package
* make graphics work with NRGBA
* how to reuse code with parametrized scenes? (e.g. effectNode)