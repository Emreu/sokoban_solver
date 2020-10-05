extends Node2D


# Declare member variables here. Examples:
# var a = 2
# var b = "text"
var tiles
var tile_size = 8


# Called when the node enters the scene tree for the first time.
func _ready():
	tiles = Image.new()
	if tiles.load("assets/preview.svg") != OK:
		print("Can't load preview tiles image")
	tiles.convert(Image.FORMAT_RGB8)
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
	tex.create_from_image(img, 0)
	return tex


#func _process(delta):
#	var tex = makeTexture(128,128)
#	$TextureRect.texture = tex

func _on_Button_pressed():
	var tex = makeTexture(128,128)
	$TextureRect.texture = tex
