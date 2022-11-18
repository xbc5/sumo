package model

func ToTag(name string) Tag {
	return Tag{
		Name: name,
	}
}

func ToTags(names []string) []Tag {
	result := []Tag{}
	for _, n := range names {
		result = append(result, Tag{Name: n})
	}
	return result
}

func ToPattern(name string, desc string, pattern string, tags []string) Pattern {
	return Pattern{
		Name:        name,
		Description: desc,
		Pattern:     pattern,
		Tags:        ToTags(tags),
	}
}
