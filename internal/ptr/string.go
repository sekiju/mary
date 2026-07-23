package ptr

func String(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
