package config

func firstNonNil(maybeNil ...interface{}) interface{} {
	for _, e := range maybeNil {
		if e != nil {
			return e
		}
	}
	return nil
}
