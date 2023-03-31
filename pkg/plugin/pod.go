package plugin

type rawPod struct {
	items      *appItems
	hasChanged *bool
}

func (r *rawPod) Add(name string, value interface{}) {
	list, name := splitasPath(name)
	// fmt.Println(">>", r.items)
	// fmt.Println("!>", list)

	path := append([]string{"spec"}, list...)

	// fmt.Println(".>", path)

	r.items.add(name, value, path)
	*r.hasChanged = true
}
