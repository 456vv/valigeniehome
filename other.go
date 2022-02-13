package valigeniehome

var HomeContextKey = &contextKey{"Home"}


type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "valigeniehome context value " + k.name }
