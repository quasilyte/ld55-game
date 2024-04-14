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

TODO:

* crt filter?
* menu bg
* "enemy in arc" condition
* instruction progression? (don't unlock all at level 1)
* weapon/vessel balancing
* vessel explosion effect (animation)

TODO (polishing):

* credits screen (you get there after winning)
* sound level balancing (sfx)
* optimize sfx (if needed)

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