extends ItemList

func Load(initial):
	clear()
	add_item("-- Start --")
	var last = 1
	for dir in initial:
		add_item(String(last) + ": " + Common.DirNames[dir])
		last+=1
	
func ClearUpto(index :int):
	for i in range(get_item_count()-1, index, -1):
		remove_item(i)
	
func Append(dir):
	var last = get_item_count()
	add_item(String(last) + ": " + Common.DirNames[dir])
	select(last)
	ensure_current_is_visible()
	
func Select(index: int):
	select(index)
	ensure_current_is_visible()
