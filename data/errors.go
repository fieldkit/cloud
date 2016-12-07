package data

type Errors map[string][]string

func NewErrors() Errors {
	return make(map[string][]string)
}

func (e Errors) Error(key, message string) {
	e[key] = append(e[key], message)
}
