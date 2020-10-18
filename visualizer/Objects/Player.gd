extends KinematicBody2D

export var DirectControl = true
var tile_size = 16
onready var ray = $RayCast2D


var inputs = {
	"ui_right": Vector2.RIGHT,
	"ui_left": Vector2.LEFT,
	"ui_up": Vector2.UP,
	"ui_down": Vector2.DOWN
}

func _unhandled_input(event):
	if not DirectControl:
		return
	for dir in inputs.keys():
		if event.is_action_pressed(dir):
			Move(inputs[dir])

func Move(dir):
	var dest = dir * tile_size
	ray.cast_to = dest
	ray.force_raycast_update()
	if !ray.is_colliding():
		position += dest
		return true
	else:
		var collider = ray.get_collider()
		if collider.is_in_group("box"):
			if collider.move(dir):
				position += dest
				return true
	return false
