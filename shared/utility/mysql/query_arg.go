package mysql

type QueryArg struct {
	Query string
	Arg   interface{}
}

func NewQueryArg(query string, arg interface{}) *QueryArg {
	return &QueryArg{
		Query: query,
		Arg:   arg,
	}
}

type QueryArgs struct {
	Query string
	Args  []interface{}
}

func NewQueryArgs(query string, args []interface{}) *QueryArgs {
	return &QueryArgs{
		Query: query,
		Args:  args,
	}
}
