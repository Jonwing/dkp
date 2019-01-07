package purge

type Op string

type InfoProvider interface {
	StringField(f string) string
	IntField(f string) int64
}

const (
	EQ Op = "="
	NE Op = "!="
	LTE Op = "<="
	GTE Op = ">="
	LT Op = "<"
	GT Op = ">"
)

type Compare func(one, another interface{}) bool

func EqInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt == SecondInt
}

func NeInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt != SecondInt
}

func LteInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt <= SecondInt
}

func GteInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt >= SecondInt
}

func LtInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt < SecondInt
}

func GtInt64(first, second interface{}) bool {
	FirstInt := first.(int64)
	SecondInt := second.(int64)
	return FirstInt > SecondInt
}

func EqString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt == SecondInt
}

func NeString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt != SecondInt
}

func LteString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt <= SecondInt
}

func GteString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt >= SecondInt
}

func LtString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt < SecondInt
}

func GtString(first, second interface{}) bool {
	FirstInt := first.(string)
	SecondInt := second.(string)
	return FirstInt > SecondInt
}


var int64Comparator = map[Op]Compare{
	EQ: EqInt64,
	NE: NeInt64,
	LTE: LteInt64,
	GTE: GteInt64,
	LT: LtInt64,
	GT: GtInt64,
}

var stringComparator = map[Op]Compare{
	EQ: EqString,
	NE: NeString,
	LTE: LteString,
	GTE: GteString,
	LT: LtString,
	GT: GtString,
}

type Filter struct {
	Source string

	Field string

	Comparator Op

	// value is parsed from string, just store the original value
	Value string
}

// parseFilter uses filterPtn to parse -f argument into Filter instance
func parseFilter(s string) (f Filter, err error) {
	m := filterPtn.FindStringSubmatch(s)
	if m == nil {
		return f, Mismatched
	}

	f = Filter{Source: s}
	for i, name := range filterPtn.SubexpNames() {
		if i != 0 && name != "" && m[i] != "" {
			switch name {
			case "field":
				f.Field = m[i]
			case "op":
				f.Comparator = Op(m[i])
			case "value":
				f.Value = m[i]
			}
		}
	}
	return
}
