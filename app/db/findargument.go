package db

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/ttacon/glog"
)

// FindArguments represents a collection of FindArgument's
type FindArguments []FindArgument

// BuildQueryArgs, given an input query and limit, returns the final
// SQL query and arguments to pass to dbmap.
func (fas FindArguments) BuildQueryArgs(inputQuery string, limit int, allowedCols map[string]interface{}) (string, map[string]interface{}) {
	var query strings.Builder
	query.WriteString(inputQuery)

	args := map[string]interface{}{
		"limit": limit,
	}

	for i, f := range fas {
		if i == 0 {
			query.WriteString(` WHERE `)
		} else {
			query.WriteString(` AND `)
		}
		query.WriteString(f.SQL(allowedCols))
		for k, v := range f.Args() {
			args[k] = v
		}
	}
	query.WriteString(` ORDER BY time DESC LIMIT :limit`)

	glog.Infoln("BuildQueryArgs:", query.String(), args)
	return query.String(), args
}

// FindArgument represents a field name, operation, and value
type FindArgument struct {
	FieldName string
	Operation string
	Value     string
	prefix    string
}

// NewFindArgument creates a FindArgument object
func NewFindArgument(name, op, val string) FindArgument {
	return FindArgument{
		FieldName: name,
		Operation: op,
		Value:     val,
		prefix:    randomString(8),
	}
}

// FindArgumentsFromQuery builds a []FindArgument and returns the count
// from a request.URL.Query() object
func FindArgumentsFromQuery(query url.Values) ([]FindArgument, int) {
	count := 10

	var findArgs []FindArgument
	for k, vals := range query {
		v := vals[0]
		if k == "count" {
			if c, err := strconv.Atoi(v); err == nil {
				count = c
			}
		} else if strings.HasPrefix(k, "find[") {
			findParts := strings.Split(strings.TrimPrefix(k, "find["), "[")
			if len(findParts) < 1 {
				continue
			}
			argName := strings.TrimSuffix(findParts[0], "]")
			argOp := "$eq"
			if len(findParts) >= 2 {
				argOp = strings.TrimSuffix(findParts[1], "]")
			}
			glog.Infoln("newFindArg: ", argName, argOp, v)
			findArgs = append(findArgs, NewFindArgument(argName, argOp, v))
		}
	}
	return findArgs, count
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

// SQL returns the SQL form string `:key (oper) :value`, and replaces the
// field name with an alias in allowedCols if one exists
func (f FindArgument) SQL(allowedCols map[string]interface{}) string {
	name := f.FieldName
	if alias, ok := allowedCols[f.FieldName]; !ok {
		glog.Fatalf("Field name %s is not in allowedCols: %s", f.FieldName, allowedCols)
		return ""
	} else if alias != nil {
		name = alias.(string)
	}

	if f.Operation == "$gte" {
		return name + ` >= :` + f.prefix + `Value`
	} else if f.Operation == "$lte" {
		return name + ` <= :` + f.prefix + `Value`
	}
	return name + ` == :` + f.prefix + `Value`
}

// Args returns the map for the arguments used as parameters for SQL
func (f FindArgument) Args() map[string]interface{} {
	return map[string]interface{}{
		f.prefix + "Value": tryInt(f.Value),
	}
}

func tryInt(s string) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return s
}
