extends WindowDialog


# Called when the node enters the scene tree for the first time.
func _ready():
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#	pass


func _on_TimeoutSlider_value_changed(value):
	var amount = String(value) + "sec."
	if value == $TimeoutSlider.max_value:
		amount = "unlimited"
	
	$TimeoutLabel.text = "Timeout: " + amount


func _on_Button_pressed():
	hide()

func GetURL():
	var params = "?"
	if $Deadzones.pressed:
		params += "deadzones=true&"
	if $Metrics.pressed:
		params += "metrics=true&"
	if $States.pressed:
		params += "states=true&"
		params += "max_states=" + String($MaxStates.value) + "&"
	if $TimeoutSlider.value != $TimeoutSlider.max_value:
		params += "timeout=" + String($TimeoutSlider.value)
	return "http://" + $URLEdit.text + "/solve" + params
