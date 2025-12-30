package ptti

import (
	"fmt"
	"maps"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Type string

const (
	String Type = "string"
	Array  Type = "array"
	ID     Type = "id"
	Int    Type = "int"
	Bool   Type = "bool"
	Date   Type = "date"
)

type QueryOpts struct {
	Filters    bson.M
	Projection bson.M
	Sort       bson.D
	Page       int64
	Limit      int64
	Offset     int64
	Skip       int64
}

func Parse(q url.Values, specs map[string]Type) (QueryOpts, error) {
	out := QueryOpts{
		Filters: bson.M{},
	}

	// base filter
	if v := strings.TrimSpace(q.Get("id")); v != "" {
		oid, err := bson.ObjectIDFromHex(v)
		if err != nil {
			return QueryOpts{}, fmt.Errorf("invalid id: %w", err)
		}
		out.Filters["_id"] = oid
	}
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		out.Filters["status"] = v
	}

	// projection
	if v := strings.TrimSpace(q.Get("fields")); v != "" {
		proj := bson.M{}
		for f := range strings.SplitSeq(v, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			proj[f] = 1
		}
		if len(proj) > 0 {
			out.Projection = proj
		}
	}

	// pagination
	if v := strings.TrimSpace(q.Get("page")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 1 {
			return QueryOpts{}, fmt.Errorf("invalid page")
		}
		out.Page = n
	}

	if v := strings.TrimSpace(q.Get("limit")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 1 {
			return QueryOpts{}, fmt.Errorf("invalid limit")
		}
		out.Limit = n
	}
	if v := strings.TrimSpace(q.Get("offset")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 1 {
			return QueryOpts{}, fmt.Errorf("invalid offset")
		}
		out.Limit = n
	}

	// Compute Skip (if offset is given)
	if out.Offset > 0 {
		out.Skip = out.Offset
	} else if out.Page > 0 {
		if out.Limit < 1 {
			return QueryOpts{}, fmt.Errorf("page requires limit")
		}
		out.Skip = (out.Page - 1) * out.Limit
	}

	// sort
	sortBy := strings.TrimSpace(q.Get("sort"))
	order := strings.ToLower(strings.TrimSpace(q.Get("order")))
	if sortBy != "" {
		dir := int32(-1)
		if order == "asc" {
			dir = 1
		}
		out.Sort = bson.D{{Key: sortBy, Value: dir}}
	}

	for k, t := range specs {
		switch t {
		case Array:
			vals, ok := q[k]
			if !ok || len(vals) == 0 {
				continue
			}
			in := make([]any, 0, len(vals))
			for _, s := range vals {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				if oid, err := bson.ObjectIDFromHex(s); err == nil {
					in = append(in, oid)
				} else {
					in = append(in, s)
				}
			}
			if len(in) > 0 {
				out.Filters[k] = bson.M{"$in": in}
			}
		case ID:
			v := strings.TrimSpace(q.Get(k))
			if v == "" {
				continue
			}
			oid, err := bson.ObjectIDFromHex(v)
			if err != nil {
				return QueryOpts{}, fmt.Errorf("invalid %s ObjectID: %w", k, err)
			}
			out.Filters[k] = oid
		case Int:
			v := strings.TrimSpace(q.Get(k))
			if v == "" {
				continue
			}
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return QueryOpts{}, fmt.Errorf("invalid %s int: %w", k, err)
			}
			out.Filters[k] = n
		case Bool:
			v := strings.TrimSpace(q.Get(k))
			if v == "" {
				continue
			}
			b, err := strconv.ParseBool(v)
			if err != nil {
				return QueryOpts{}, fmt.Errorf("invalid %s bool: %w", k, err)
			}
			out.Filters[k] = b
		case Date:
			q_date := strings.TrimSpace(q.Get(k))
			if q_date != "" {
				dateFmt, err := time.Parse(time.DateOnly, q_date)
				if err != nil {
					continue
				}
				dateTimeStart := dateFmt.Format(time.RFC3339)
				date, err := time.Parse(time.RFC3339, dateTimeStart)
				if err != nil {
					continue
				}
				out.Filters[k] = date
			}
		// string is default
		default:
			v := strings.TrimSpace(q.Get(k))
			if v == "" {
				continue
			}
			out.Filters[k] = v
		}
	}
	return out, nil
}

func MongoFind(qo QueryOpts, extra bson.M) (bson.M, *options.FindOptionsBuilder) {
	filters := bson.M{}
	for _, m := range []bson.M{qo.Filters, extra} {
		maps.Copy(filters, m)
	}
	opts := options.Find()
	if len(qo.Projection) > 0 {
		opts.SetProjection(qo.Projection)
	}
	if len(qo.Sort) > 0 {
		opts.SetSort(qo.Sort)
	}
	if qo.Limit > 0 {
		opts.SetLimit(qo.Limit)
	}
	if qo.Skip > 0 {
		opts.SetSkip(qo.Skip)
	}
	return filters, opts
}
