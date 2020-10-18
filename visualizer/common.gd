class_name Common

enum Direction {UP, RIGHT, DOWN, LEFT}

const DirNames = {
	Direction.UP: "Up",
	Direction.DOWN: "Down",
	Direction.LEFT: "Left",
	Direction.RIGHT: "Right",
}

const DirVec = {
	Direction.UP: Vector2(0, -1),
	Direction.DOWN: Vector2(0, 1),
	Direction.LEFT: Vector2(-1, 0),
	Direction.RIGHT: Vector2(1, 0),
}
