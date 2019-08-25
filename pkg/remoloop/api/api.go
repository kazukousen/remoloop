package api

type Resource string

func (r Resource) String() string {
	return string(r)
}
