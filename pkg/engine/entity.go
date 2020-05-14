package engine

// Entity is an interface that aides the engine in having commmon functionality
type Entity interface {
	Draw()
	SetX(float64)
	SetY(float64)
	Update()
}
