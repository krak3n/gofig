package gofig

type tag struct {
	name string
}

func parseTag(t string) tag {
	return tag{
		name: t,
	}
}
