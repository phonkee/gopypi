package classy

/*
return MultiRoutes().
    Add(Detail.Routes(), "{{view}}_detail").
    Add(List.Routes(), "{{view}}_list").
    Get()
*/

func JoinRoutes() multiroute {
	return multiroute{
		m: map[string]Mapping{},
	}
}

type multiroute struct {
	m map[string]Mapping
}

func (m multiroute) Add(route map[string]Mapping, name string, aliases ...[]string) multiroute {
	for k, v := range route {

		m.m[k] = v.Name(name)

		for _, alias := range aliases {
			if len(alias) != 2 {
				panic("aliases must be 2 elements (from, to)")
			}

			m.m[k] = m.m[k].Rename(alias[0], alias[1])
		}
	}
	return m
}

func (m multiroute) Get() map[string]Mapping {
	return m.m
}
