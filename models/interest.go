package models

type Interest struct {
	ID   int
	Name string
}

func (i *Interest) GetScanForm() []interface{} {
	return []interface{}{
		&i.ID,
		&i.Name,
	}
}
