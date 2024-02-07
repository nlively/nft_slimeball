package main

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	GRAVITY           = 9.8 * 3  // Gravity acceleration
	TIME_STEP         = 1 / 30.0 // Time step for the simulation, in seconds
	BOUNCE_EFFICIENCY = .85      // To simulate energy loss on bounce, 1.0 for a perfect bounce, less for damping
	FLOOR_FRICTION    = 0.98

	STATE_START = 0
	STATE_SERVE = 1
	STATE_PLAY  = 2

	BALL_RESTING = 0
	BALL_MOVING  = 1
)
