extends Control

onready var Moves = $VBoxContainer/HBoxContainer/RightPanel/Moves
onready var Level = $VBoxContainer/HBoxContainer/MarginContainer/ViewportContainer/LevelView/Level
var currentMove = -1
var moves = []
var states = []
var levelPath = ""

var deadZones = []
var metrics = []

signal exit

func Init(level):
	currentMove = -1
	moves = []
	states = []
	deadZones = []
	metrics = []
	levelPath = level["path"]
	Level.Load(level)
	for node in get_tree().get_nodes_in_group("temp"):
		node.queue_free()

func _input(event):
	if event.is_action_pressed("ui_up"):
		DoMove(Common.Direction.UP)
	if event.is_action_pressed("ui_down"):
		DoMove(Common.Direction.DOWN)
	if event.is_action_pressed("ui_left"):
		DoMove(Common.Direction.LEFT)
	if event.is_action_pressed("ui_right"):
		DoMove(Common.Direction.RIGHT)
	if event.is_action_pressed("ui_undo") && currentMove > 0:
		GotoMove(currentMove-1)

func _on_Reset_pressed():
	moves.clear()
	states.clear()
	Moves.ClearUpto(0)
	currentMove = -1

func DoMove(dir):
	# check if current move is valid in
	var vec = Common.DirVec[dir]
	var done = Level.TryMove(vec)
	if not done:
		return
	# If current move is not last - clean all subsequent moves
	if currentMove+1 != moves.size():
		moves.resize(currentMove)
		Moves.ClearUpto(currentMove)
	# save move
	moves.append(dir)
	Moves.Append(dir)
	currentMove = moves.size()
	# save state
	var state = Level.GetState()
	states.append(state)
	$VBoxContainer/HBoxContainer/RightPanel/CurMove.text = String(currentMove)

func GotoMove(index: int):
	print(index)
	currentMove = index
	# TODO: FastForward states
	Level.SetState(states[index])
	$VBoxContainer/HBoxContainer/RightPanel/CurMove.text = String(currentMove)
	Moves.Select(currentMove)


func _on_Debug_pressed():
	var output = []
	var exitCode = OS.execute("/home/emreu/otus/sokoban_solver/solver/solver.out", ["-f", levelPath, "-debug"], true, output)
	print("exit code:" + String(exitCode))
	print(output)
	var res = JSON.parse(output[0])
	if res.error != OK:
		print("Can't parse json:" + res.error_string)
		return
		
	if res.result.has("dead_zones"):
		deadZones.clear()
		for pos in res.result["dead_zones"]:
			deadZones.append(Vector2(pos["X"], pos["Y"]))
	if res.result.has("metrics"):
		var n = 0
		metrics.clear()
		for data in res.result["metrics"]:
			var goalMetric = {}
			for p in data:
				var coords = p.split_floats(",")
				goalMetric[Vector2(coords[0], coords[1])] = data[p]
			metrics.append(goalMetric)
			var btn = Button.new()
			btn.text = String(n)
			btn.add_to_group("temp")
			btn.connect("pressed", self, "_on_MetricBtn_pressed", [n])
			$VBoxContainer/ButtonsBar.add_child(btn)
			n+=1

func _on_Back_pressed():
	emit_signal("exit")

func _on_MetricBtn_pressed(index):
	print("Show metric #" + String(index))
	Level.ShowMetrics(metrics[index])

func _on_DeadZones_toggled(button_pressed):
	if button_pressed:
		Level.ShowDeadzones(deadZones)
	else:
		Level.ShowDeadzones([])

func _on_NoMetrics_pressed():
	Level.ShowMetrics({})


func _on_Tree_pressed():
	var output = []
	var exitCode = OS.execute("/home/emreu/otus/sokoban_solver/solver/solver.out",
		["-f", levelPath, "-tree", "-max", "1000", "-timeout", "10s"], true, output)
	print("exit code:" + String(exitCode))
	print(len(output))
	
	var res = JSON.parse(output[0])
	if res.error != OK:
		print("Can't parse json:" + res.error_string)
		return
	
	DrawTree(res.result)
	$StateTreeWindow.show()
	
func DrawTree(states):
	var parents = {}
	for state in states:
		var parent = null
		if state["parent"] > 0:
			parent = parents[state.parent]
		var item = $StateTreeWindow/Tree.create_item(parent)
		item.collapsed = true
		item.set_text(0, "state #" + String(state["id"]))
		item.set_tooltip(0, String(state["hash"]))
		item.set_text(1, String(state["metric"]))
		var boxes = []
		for box in state["boxes"]:
			boxes.append(Vector2(box["X"], box["Y"]))
		item.set_metadata(0, boxes)
		var domain = []
		for pos in state["domain"]:
			domain.append(Vector2(pos["X"], pos["Y"]))
		item.set_metadata(1, domain)
		parents[state["id"]] = item

func _on_Tree_item_selected():
	var item = $StateTreeWindow/Tree.get_selected()
	if not item:
		return
	var boxes = item.get_metadata(0)
	var domain = item.get_metadata(1)
	
	Level.SetState({"boxes": boxes})
	Level.ShowDeadzones(domain)
	
