extends Control

onready var Moves = $VBoxContainer/HBoxContainer/RightPanel/Moves
onready var Level = $VBoxContainer/HBoxContainer/MarginContainer/ViewportContainer/LevelView
var currentMove = -1
var moves = []
var states = []

func _input(event):
	if event.is_action_pressed("ui_up"):
		AddMove(Common.Direction.UP)
	if event.is_action_pressed("ui_down"):
		AddMove(Common.Direction.DOWN)
	if event.is_action_pressed("ui_left"):
		AddMove(Common.Direction.LEFT)
	if event.is_action_pressed("ui_right"):
		AddMove(Common.Direction.RIGHT)
	if event.is_action_pressed("ui_undo") && currentMove > 0:
		GotoMove(currentMove-1)

func _on_Reset_pressed():
	moves.clear()
	Moves.ClearUpto(0)
	currentMove = -1

func AddMove(dir):
	# check if current move is valid in 
	# If current move is not last - clean all subsequent moves
	if currentMove+1 != moves.size():
		moves.resize(currentMove)
		Moves.ClearUpto(currentMove)
	
	moves.append(dir)
	Moves.Append(dir)
	currentMove = moves.size()
	$VBoxContainer/HBoxContainer/RightPanel/CurMove.text = String(currentMove)

func GotoMove(index: int):
	print(index)
	currentMove = index
	$VBoxContainer/HBoxContainer/RightPanel/CurMove.text = String(currentMove)
	Moves.Select(currentMove)
