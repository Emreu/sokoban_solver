extends Node2D

var box = preload("res://Objects/Box.tscn")
var goal = preload("res://Objects/Goal.tscn")
var wall = preload("res://Objects/Wall.tscn")

var zone = preload("res://Objects/Deadzone.tscn")

const tileSize = Vector2(16, 16)

signal win

func _ready():
	var level = {
		"walls": [Vector2(2,2)],
		"goals": [Vector2(3,3)],
		"boxes": [Vector2(4,4)],
		"start": Vector2(0,0),
	}
	Load(level)

func Load(level):
	# clean
	for w in $Walls.get_children():
		w.queue_free()
	for b in $Boxes.get_children():
		b.queue_free()
	for g in $Goals.get_children():
		g.queue_free()
	for z in $Deadzones.get_children():
		z.queue_free()
	for m in $Metrics.get_children():
		m.queue_free()
	
	# start position
	$Player.position = tileSize * level["start"]
	
	# add walls
	for pos in level["walls"]:
		var w = wall.instance()
		w.position = tileSize * pos
		$Walls.add_child(w, true)
	# add goals
	for pos in level["goals"]:
		var g = goal.instance()
		g.position = tileSize * pos
		$Goals.add_child(g, true)
	# add boxes
	for pos in level["boxes"]:
		var b = box.instance()
		b.position = tileSize * pos
		$Boxes.add_child(b, true)
	# position camera
	$Camera.zoom = Vector2(0.5,0.5)

func TryMove(dir):
	var moved = $Player.Move(dir)
	if moved:
		call_deferred("checkWin")
	return moved
	
func GetState():
	var boxes = []
	for box in $Boxes.get_children():
		boxes.append(box.position / tileSize)
	return {
		"boxes": boxes,
		"player": $Player.position / tileSize,
	}
	
func SetState(state):
	for box in state["boxes"]:
		pass
	pass
	
func ShowDeadzones(positions):
	for z in $Deadzones.get_children():
		z.queue_free()
	for pos in positions:
		var z = zone.instance()
		z.position = tileSize * pos
		$Deadzones.add_child(z, true)
	
func checkWin():
	var boxPositions = []
	for box in $Boxes.get_children():
		boxPositions.append(box.position)
	for goal in $Goals.get_children():
		if !boxPositions.has(goal.position):
			return
	emit_signal("win")
	
func ShowMetrics(metrics):
	for m in $Metrics.get_children():
		m.queue_free()
	for pos in metrics:
		var m = Label.new()
		m.text = String(metrics[pos])
		m.rect_position = tileSize * pos
		$Metrics.add_child(m, true)
	
func _on_Level_win():
	print("Win!")
