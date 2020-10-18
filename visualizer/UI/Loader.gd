extends Control

var levels = []
var currentIndex = -1

# Called when the node enters the scene tree for the first time.
func _ready():
	LoadMap("/home/emreu/otus/sokoban_solver/solver/test_1.txt")
	LoadMap("/home/emreu/otus/sokoban_solver/solver/test_2.txt")
	LoadMap("/home/emreu/otus/sokoban_solver/solver/test_3.txt")
	pass # Replace with function body.

func LoadMap(path):
	var file = File.new()
	file.open(path, file.READ)
	var content = file.get_as_text()
	file.close()
	
	var level = parseLevel(content)
	if not level:
		print("level not loaded!")
		return
	
	level["path"] = path
		
	var name = path.get_file()
	if not name:
		name = "Level " + String(levels.size())
		
	levels.append(level)
	$LoaderUI/LevelsList.add_item(name)
	
func parseLevel(chars):
	var walls = []
	var goals = []
	var boxes = []
	var start = Vector2()
	var pos = Vector2(-1,0)
	for c in chars:
		pos.x += 1
		match c:
			"#":
				walls.append(pos)
			"@":
				start = pos
			"+":
				start = pos
				goals.append(pos)
			"$":
				boxes.append(pos)
			"*":
				boxes.append(pos)
				goals.append(pos)
			".":
				goals.append(pos)
			"\n":
				pos.x = -1
				pos.y += 1

	return {
		"walls": walls,
		"goals": goals,
		"boxes": boxes,
		"start": start,
	}

func _on_MapLoadButton_pressed():
	$LoaderUI/FileDialog.show()

func _on_FileDialog_file_selected(path):
	LoadMap(path)


func _on_LevelsList_item_selected(index):
	currentIndex = index
	var level = levels[index]
	$LoaderUI/Preview.ShowLevel(level)

func _on_PlayButton_pressed():
	if currentIndex < 0:
		return
	$LoaderUI.hide()
	var level = levels[currentIndex]
	$Viewer.Init(level)
	$Viewer.show()

func _on_Viewer_exit():
	$Viewer.hide()
	$LoaderUI.show()
