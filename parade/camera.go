package parade

type Camera interface {
	Position() (float64, float64, float64)
	SetPosition(float64, float64, float64)
	Direction() (float64, float64, float64)
	SetDirection(float64, float64, float64)
}
