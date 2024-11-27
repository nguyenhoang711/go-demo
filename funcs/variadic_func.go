package fuct

func Join(del string, values []string) string {
	var line string
	for i, v := range values {
		line = line + v
		if i != len(values) - 1 {
			line = line + del
		}
	}
	return line
}