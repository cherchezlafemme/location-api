// Copyright 2023 The Infratographer Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.infratographer.com/location-api/internal/ent/generated/location"
	"go.infratographer.com/x/gidx"
)

// Common entgql types.
type (
	Cursor         = entgql.Cursor[gidx.PrefixedID]
	PageInfo       = entgql.PageInfo[gidx.PrefixedID]
	OrderDirection = entgql.OrderDirection
)

func orderFunc(o OrderDirection, field string) func(*sql.Selector) {
	if o == entgql.OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// LocationEdge is the edge representation of Location.
type LocationEdge struct {
	Node   *Location `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// LocationConnection is the connection containing edges to Location.
type LocationConnection struct {
	Edges      []*LocationEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

func (c *LocationConnection) build(nodes []*Location, pager *locationPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Location
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Location {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Location {
			return nodes[i]
		}
	}
	c.Edges = make([]*LocationEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &LocationEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// LocationPaginateOption enables pagination customization.
type LocationPaginateOption func(*locationPager) error

// WithLocationOrder configures pagination ordering.
func WithLocationOrder(order *LocationOrder) LocationPaginateOption {
	if order == nil {
		order = DefaultLocationOrder
	}
	o := *order
	return func(pager *locationPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultLocationOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithLocationFilter configures pagination filter.
func WithLocationFilter(filter func(*LocationQuery) (*LocationQuery, error)) LocationPaginateOption {
	return func(pager *locationPager) error {
		if filter == nil {
			return errors.New("LocationQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type locationPager struct {
	reverse bool
	order   *LocationOrder
	filter  func(*LocationQuery) (*LocationQuery, error)
}

func newLocationPager(opts []LocationPaginateOption, reverse bool) (*locationPager, error) {
	pager := &locationPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultLocationOrder
	}
	return pager, nil
}

func (p *locationPager) applyFilter(query *LocationQuery) (*LocationQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *locationPager) toCursor(l *Location) Cursor {
	return p.order.Field.toCursor(l)
}

func (p *locationPager) applyCursors(query *LocationQuery, after, before *Cursor) (*LocationQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultLocationOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *locationPager) applyOrder(query *LocationQuery) *LocationQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultLocationOrder.Field {
		query = query.Order(DefaultLocationOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *locationPager) orderExpr(query *LocationQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultLocationOrder.Field {
			b.Comma().Ident(DefaultLocationOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Location.
func (l *LocationQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...LocationPaginateOption,
) (*LocationConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newLocationPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if l, err = pager.applyFilter(l); err != nil {
		return nil, err
	}
	conn := &LocationConnection{Edges: []*LocationEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = l.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if l, err = pager.applyCursors(l, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		l.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := l.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	l = pager.applyOrder(l)
	nodes, err := l.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

var (
	// LocationOrderFieldCreatedAt orders Location by created_at.
	LocationOrderFieldCreatedAt = &LocationOrderField{
		Value: func(l *Location) (ent.Value, error) {
			return l.CreatedAt, nil
		},
		column: location.FieldCreatedAt,
		toTerm: location.ByCreatedAt,
		toCursor: func(l *Location) Cursor {
			return Cursor{
				ID:    l.ID,
				Value: l.CreatedAt,
			}
		},
	}
	// LocationOrderFieldUpdatedAt orders Location by updated_at.
	LocationOrderFieldUpdatedAt = &LocationOrderField{
		Value: func(l *Location) (ent.Value, error) {
			return l.UpdatedAt, nil
		},
		column: location.FieldUpdatedAt,
		toTerm: location.ByUpdatedAt,
		toCursor: func(l *Location) Cursor {
			return Cursor{
				ID:    l.ID,
				Value: l.UpdatedAt,
			}
		},
	}
	// LocationOrderFieldName orders Location by name.
	LocationOrderFieldName = &LocationOrderField{
		Value: func(l *Location) (ent.Value, error) {
			return l.Name, nil
		},
		column: location.FieldName,
		toTerm: location.ByName,
		toCursor: func(l *Location) Cursor {
			return Cursor{
				ID:    l.ID,
				Value: l.Name,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f LocationOrderField) String() string {
	var str string
	switch f.column {
	case LocationOrderFieldCreatedAt.column:
		str = "CREATED_AT"
	case LocationOrderFieldUpdatedAt.column:
		str = "UPDATED_AT"
	case LocationOrderFieldName.column:
		str = "NAME"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f LocationOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *LocationOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("LocationOrderField %T must be a string", v)
	}
	switch str {
	case "CREATED_AT":
		*f = *LocationOrderFieldCreatedAt
	case "UPDATED_AT":
		*f = *LocationOrderFieldUpdatedAt
	case "NAME":
		*f = *LocationOrderFieldName
	default:
		return fmt.Errorf("%s is not a valid LocationOrderField", str)
	}
	return nil
}

// LocationOrderField defines the ordering field of Location.
type LocationOrderField struct {
	// Value extracts the ordering value from the given Location.
	Value    func(*Location) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) location.OrderOption
	toCursor func(*Location) Cursor
}

// LocationOrder defines the ordering of Location.
type LocationOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *LocationOrderField `json:"field"`
}

// DefaultLocationOrder is the default ordering of Location.
var DefaultLocationOrder = &LocationOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &LocationOrderField{
		Value: func(l *Location) (ent.Value, error) {
			return l.ID, nil
		},
		column: location.FieldID,
		toTerm: location.ByID,
		toCursor: func(l *Location) Cursor {
			return Cursor{ID: l.ID}
		},
	},
}

// ToEdge converts Location into LocationEdge.
func (l *Location) ToEdge(order *LocationOrder) *LocationEdge {
	if order == nil {
		order = DefaultLocationOrder
	}
	return &LocationEdge{
		Node:   l,
		Cursor: order.Field.toCursor(l),
	}
}
