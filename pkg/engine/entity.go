package engine

// Entity is an interface that aides the engine in having commmon functionality
type Entity interface {
	Draw()
	Update()
}
