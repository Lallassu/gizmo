# Gizmo - Go 2D pixel game
This is a game under heavy development and commits might not be explained etc. Use/fork at your own risk!

Developed using faiface/pixel 2D library (https://github.com/faiface/pixel).

## Run from code
The Pixel commit below is the latest (currently) and contains a patch I've added for supporting export of `vPosition` between vertex and fragment shader and is needed by this game.
```
go mod init gizmo
go get github.com/faiface/pixel@e51d4a6676fa48c83b5ea703cb5b044e2967cb83
go run .
``` 

## Game controls 

The main menu is not ready yet, you can use `ESC` to start the game & arrows for navigation.
- b -> pick weapon
- v -> drop weapon
- q -> quit
- s -> zoom out 
- w -> zoom in
- Left Shift -> shoot

## Screenshots
![](https://raw.github.com/lallassu/gizmo/master/preview.png)

## Videos
[![](https://raw.github.com/lallassu/gizmo/master/videopreview.png)](https://youtu.be/6zcQvsf4R4Q)

## License
MIT
