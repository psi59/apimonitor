package rsdb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsvalid"
)

type queryOperator string

const (
	AndQueryOperator queryOperator = "AND"
	OrQueryOperator  queryOperator = "OR"
)

type QueryBuilder interface {
	Or(q Query) Query
	And(q Query) Query
}

type Query interface {
	QueryBuilder
	Where() string
	Values() []interface{}
}

type queryImpl struct {
	where  string
	values []interface{}
}

func (query *queryImpl) combineQueryWithOperator(q Query, operator queryOperator) Query {
	if w := query.Where(); rsvalid.IsZero(w) {
		query.where = q.Where()
	} else {
		query.where = fmt.Sprintf("(%s) %s (%s)", query.Where(), operator, q.Where())
	}
	query.values = append(query.values, q.Values()...)
	return query
}

func (query *queryImpl) Or(q Query) Query {
	return query.combineQueryWithOperator(q, OrQueryOperator)
}

func (query *queryImpl) And(q Query) Query {
	return query.combineQueryWithOperator(q, AndQueryOperator)
}

func (query *queryImpl) Where() string {
	return query.where
}

func (query *queryImpl) Values() []interface{} {
	return query.values
}

func NewQuery(w string, values ...interface{}) (Query, error) {
	if rsvalid.IsZero(w) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "invalid wheres")
	}
	if strings.Count(w, "?") != len(values) {
		return nil, errors.Wrap(rserrors.ErrInvalidParameter, "invalid values")
	}
	return &queryImpl{
		where:  w,
		values: values,
	}, nil
}

func NewEmptyQuery() Query {
	return &queryImpl{
		values: make([]interface{}, 0),
	}
}
