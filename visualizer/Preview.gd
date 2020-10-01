extends Node2D


# Declare member variables here. Examples:
# var a = 2
# var b = "text"
var tiles
var tile_size = 8


# Called when the node enters the scene tree for the first time.
func _ready():
	tiles = Image.new()
	print(tiles.load("assets/preview.svg"))
	tiles.convert(Image.FORMAT_RGB8)
	print(tiles.data.format)
	var tex = makeTexture(128,128)
	$TextureRect.texture = tex
	

func makeTexture(w, h):
	var img = Image.new()
	img.create(w, h, false, Image.FORMAT_RGB8)
	for y in range(h):
		for x in range(w):
			var i = randi() % 8
			img.blit_rect(tiles, Rect2(i*tile_size, 0, 8, 8), Vector2(x*tile_size, y*tile_size))
	var tex = ImageTexture.new()
	print(tex.create_from_image(img))
	return tex

# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#	pass

func _process(delta):
	var tex = makeTexture(128,128)
	$TextureRect.texture = tex

func _on_Button_pressed():
	var tex = makeTexture(128,128)
	$TextureRect.texture = tex
