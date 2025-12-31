package entities

func RetrieveAll() []any {
	return []any{
		&User{},
		&Post{},
		&Comment{},
	}
}
