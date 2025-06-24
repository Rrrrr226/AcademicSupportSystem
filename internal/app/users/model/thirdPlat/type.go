package thirdPlat

//go:generate go run golang.org/x/tools/cmd/stringer -type=Type -linecomment
type Type int

const (
	NotExists Type = iota // NotExists
	HDUHelp               // HDUHelp
)

var (
	List = []Type{
		NotExists,
		HDUHelp,
	}
	Map = map[string]Type{}
)

func init() {
	for _, t := range List {
		Map[t.String()] = t
	}
}
func FromString(t string) Type {
	return Map[t]
}
