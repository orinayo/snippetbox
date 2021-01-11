package forms

type errors map[string][]string

func (err errors) Add(field, message string) {
	err[field] = append(err[field], message)
}

func (err errors) Get(field string) string {
	errSlice := err[field]
	if len(errSlice) == 0 {
		return ""
	}

	return errSlice[0]
}