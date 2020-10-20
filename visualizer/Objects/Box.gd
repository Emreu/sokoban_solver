extends KinematicBody2D

var tile_size = 64
onready var ray = $RayCast2D

func move(dir):
	var dest = dir * tile_size
	ray.cast_to = dest
	ray.force_raycast_update()
	if !ray.is_colliding():
		position += dest
		return true
	return false
