package template

type Renderer interface {
	Render(user User, currency Currency) (string, error)
}
