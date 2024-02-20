package main

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	SCREEN_WIDTH_FEET  = 300.0
	SCREEN_HEIGHT_FEET = SCREEN_WIDTH_FEET * (SCREEN_HEIGHT / SCREEN_WIDTH)

	NET_HEIGHT = SCREEN_HEIGHT / 4
	NET_WIDTH  = 15

	GRAVITY           = 9.8      // Gravity acceleration
	TIME_STEP         = 1 / 20.0 // Time step for the simulation, in seconds
	BOUNCE_EFFICIENCY = 0.85     // To simulate energy loss on bounce, 1.0 for a perfect bounce, less for damping
	FLOOR_FRICTION    = 0.98

	STATE_START = 0
	STATE_SERVE = 1
	STATE_PLAY  = 2

	BALL_RESTING = 0
	BALL_MOVING  = 1
)
