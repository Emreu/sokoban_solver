extends TextureRect

var tiles = preload("res://assets/preview.svg")
const tile_size = 8
const offset = {
	"empty": 0,
	"wall": 1,
	"player": 2,
	"goal": 3,
	"player_on_goal": 4,
	"box": 5,
	"box_on_goal": 6,
}

func _ready():
	tiles.convert(Image.FORMAT_RGB8)
	
func tileAt(pos, level):
	if pos in level["walls"]:
		return "wall"
		
	var isBox = pos in level["boxes"]
	var isGoal = pos in level["goals"]
	var isPlayer = pos == level["start"]
	
	if isBox:
		if isGoal:
			return "box_on_goal"
		else:
			return "box"
			
	if isPlayer:
		if isGoal:
			return "player_on_goal"
		else:
			return "player"
			
	if isGoal:
		return "goal"
			
	return "empty"

func makeTexture(w, h, level):
	var img = Image.new()
	var size = max(w,h)
	img.create(size*tile_size, size*tile_size, false, Image.FORMAT_RGB8)
	var shift = Vector2(0,0)
	if w > h:
		var delta = (w-h)/2 * tile_size
		shift = Vector2(0,delta)
	elif h > w:
		var delta = (h-w)/2 * tile_size
		shift = Vector2(delta,0)
	
	for y in range(h):
		for x in range(w):
			var pos = Vector2(x, y)
			var tile = tileAt(pos, level)
			img.blit_rect(tiles, Rect2(offset[tile]*tile_size, 0, 8, 8), shift+Vector2(pos.x*tile_size, pos.y*tile_size))
			
	var tex = ImageTexture.new()
	tex.create_from_image(img, 0)
	
	return tex

func ShowLevel(level):
	if level.has("preview"):
		texture = level["preview"]
		return
		
	var w = 0
	var h = 0
	for pos in level["walls"]:
		if pos.x > w:
			w = pos.x
		if pos.y > h:
			h = pos.y
	var tex = makeTexture(w+1,h+1,level)
	level["preview"] = tex
	texture = tex

