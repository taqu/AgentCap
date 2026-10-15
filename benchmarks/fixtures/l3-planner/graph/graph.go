package graph
type Node struct { Name string; Deps []string }
type Graph map[string]Node
func(g Graph)Order()([]string,error){var names []string;for n:=range g{names=append(names,n)};return names,nil}
