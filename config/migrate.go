package config

type migration struct {
	from, to int
	apply    func(raw map[string]any) error
}

var migrations = []migration{
	// Empty at launch — add entries here for future JSON schema changes.
	// Example: {from: 1, to: 2, apply: migrateV1ToV2}
}

func RunMigrations(raw map[string]any) (changed bool, err error) {
	version, _ := raw["version"].(float64)
	if version <= 0 {
		version = float64(CurrentConfigVersion)
	}

	for _, m := range migrations {
		if version >= float64(m.to) {
			continue
		}
		if version != float64(m.from) {
			continue
		}
		if err := m.apply(raw); err != nil {
			return changed, err
		}
		raw["version"] = float64(m.to)
		version = float64(m.to)
		changed = true
	}

	return changed, nil
}
