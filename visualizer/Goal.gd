extends Node2D

var occupied = false

onready var sprite_on = $sprite_on
onready var sprite_off = $sprite_off


func _on_Detector_body_entered(body):
	if body.is_in_group("box"):
		occupied = true
		sprite_on.show()
		sprite_off.hide()

func _on_Detector_body_exited(body):
	if body.is_in_group("box"):
		occupied = false
		sprite_off.show()
		sprite_on.hide()
