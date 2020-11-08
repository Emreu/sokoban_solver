extends Control

onready var Moves = $VBoxContainer/HBoxContainer/RightPanel/Moves
onready var Level = $VBoxContainer/HBoxContainer/MarginContainer/ViewportContainer/LevelView/Level
var currentState = 0
var states = []
var level = null

var debug_deadzones = []
var debug_metrics = []
var debug_hlock = []
var debug_vlock = []

signal exit

func Init(newLevel):
	level = newLevel
	currentState = 0
	states = []
	Moves.ClearUpto(0)
	
	debug_deadzones = []
	debug_metrics = []
	debug_hlock = []
	debug_vlock = []
	
	Level.Load(newLevel)
	# wait for 1 frame so old object will be deleted from level
	yield(get_tree(),"idle_frame")
	states.append(Level.GetState())
	
	for node in get_tree().get_nodes_in_group("temp"):
		node.queue_free()
	for node in get_tree().get_nodes_in_group("default_disabled"):
		node.disabled = true

func _input(event):
	if event.is_action_pressed("ui_up"):
		DoMove(Common.Direction.UP)
	if event.is_action_pressed("ui_down"):
		DoMove(Common.Direction.DOWN)
	if event.is_action_pressed("ui_left"):
		DoMove(Common.Direction.LEFT)
	if event.is_action_pressed("ui_right"):
		DoMove(Common.Direction.RIGHT)
	if event.is_action_pressed("ui_undo") && currentState > 0:
		GotoState(currentState-1)
	if event.is_action_pressed("ui_next") && currentState+1 < states.size():
		GotoState(currentState+1)

func _on_Reset_pressed():
	var initial = states[0]
	states.clear()
	states.append(initial)
	Moves.ClearUpto(0)
	currentState = 0
	yield(get_tree(),"idle_frame")
	yield(get_tree(),"physics_frame")

func DoMove(dir):
	# check if current move is valid in
	var vec = Common.DirVec[dir]
	var done = Level.TryMove(vec)
	if not done:
		return
	# If current move is not last - clean all subsequent moves
	if currentState+1 < states.size():
		states.resize(currentState)
		Moves.ClearUpto(currentState)
	# save move
	Moves.Append(dir)
	currentState += 1
	# save state
	var state = Level.GetState()
	states.append(state)

func GotoState(index: int):
	currentState = index
	Level.SetState(states[index])
	Moves.Select(currentState)

func _on_Back_pressed():
	emit_signal("exit")

func _on_MetricBtn_pressed(index):
	print("Show metric #" + String(index))
	Level.ShowMetrics(debug_metrics[index])

func _on_DeadZones_toggled(button_pressed):
	if button_pressed:
		Level.ShowDeadzones(debug_deadzones)
	else:
		Level.ShowDeadzones([])

func _on_NoMetrics_pressed():
	Level.ShowMetrics({})

func _on_Tree_pressed():
	$StateTreeWindow.popup()
	
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
		if state.has("fail"):
			item.set_text(2, String(state["fail"]))
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

func _on_Solve_pressed():
	var url = $SettingsWindow.GetURL()
	var levelASCII = level["source"]
	
	var err = $HTTPRequest.request(url, [], false, HTTPClient.METHOD_POST, levelASCII)
	if err != OK:
		print(err)

func _on_Settings_pressed():
	$SettingsWindow.popup_centered()
	pass # Replace with function body.

func _on_HTTPRequest_request_completed(result, response_code, headers, body):
	print(result, response_code)
	if result != OK:
		print("Not OK!")
		return
	
	if response_code != 200:
		print("Error: ", body.get_string_from_ascii())
		return
		
	var res = JSON.parse(body.get_string_from_ascii())
	if res.error != OK:
		print("Can't parse json:" + res.error_string)
		return
	
	# save debug data if present
	if res.result.has("dead_zones"):
		debug_deadzones.clear()
		for pos in res.result["dead_zones"]:
			debug_deadzones.append(Vector2(pos["X"], pos["Y"]))
		$VBoxContainer/ButtonsBar/DeadZones.disabled = false
	if res.result.has("horizontal_locks"):
		debug_hlock.clear()
		for pos in res.result["horizontal_locks"]:
			debug_hlock.append(Vector2(pos["X"], pos["Y"]))
		$VBoxContainer/ButtonsBar/HVLocks.disabled = false
	if res.result.has("vertical_locks"):
		debug_vlock.clear()
		for pos in res.result["vertical_locks"]:
			debug_vlock.append(Vector2(pos["X"], pos["Y"]))
		$VBoxContainer/ButtonsBar/HVLocks.disabled = false
	if res.result.has("metrics"):
		var n = 0
		debug_metrics.clear()
		for data in res.result["metrics"]:
			var goalMetric = {}
			for p in data:
				var coords = p.split_floats(",")
				goalMetric[Vector2(coords[0], coords[1])] = data[p]
			debug_metrics.append(goalMetric)
			var btn = Button.new()
			btn.text = String(n)
			btn.add_to_group("temp")
			btn.connect("pressed", self, "_on_MetricBtn_pressed", [n])
			$VBoxContainer/ButtonsBar.add_child(btn)
			n+=1
		$VBoxContainer/ButtonsBar/NoMetrics.disabled = false
	if res.result.has("states"):
		DrawTree(res.result["states"])
		$VBoxContainer/ButtonsBar/Tree.disabled = false
	
	# replay solution
	if res.result.has("solution"):
		for c in res.result["solution"]:
			match c:
				"^":
					DoMove(Common.Direction.UP)
				">":
					DoMove(Common.Direction.RIGHT)
				"V":
					DoMove(Common.Direction.DOWN)
				"<":
					DoMove(Common.Direction.LEFT)
			yield(get_tree(),"idle_frame")
			yield(get_tree(),"physics_frame")


func _on_HVLocks_toggled(button_pressed):
	if button_pressed:
		Level.ShowLocks(debug_hlock, debug_vlock)
	else:
		Level.ShowLocks([], [])
