package helper

import (
	"log"
	"maps"
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func baseQuery(q url.Values) map[string]any {
	filters := bson.M{}
	baseMap := make(map[string]any)
	if q.Has("id") {
		objectID, err := bson.ObjectIDFromHex(q.Get("id"))
		if err != nil {
			log.Printf("error decode %s", err)
		}
		filters["_id"] = objectID
	}
	if q.Has("status") {
		status := q.Get("status")
		filters["status"] = status
	}
	if q.Has("fields") {
		arrStr := strings.Split(q.Get("fields"), ",")
		baseMap["fields"] = arrStr
	}
	if q.Has("page") {
		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			log.Printf("err: %s", err)
		}
		baseMap["page"] = page
	}
	if q.Has("offset") {
		offset, err := strconv.Atoi(q.Get("offset"))
		if err != nil {
			log.Printf("err: %s", err)
		}
		baseMap["offset"] = offset
	}
	if q.Has("sort") && q.Has("order") {
		sortBy := q.Get("sort")
		var orderBy int
		switch strings.ToLower(q.Get("order")) {
		case "asc":
			orderBy = 1
		default:
			orderBy = -1
		}
		sortMap := map[string]int{sortBy: orderBy}

		baseMap["sort"] = sortMap
	}

	baseMap["filters"] = filters

	return baseMap
}

// ParseQuery is a parser for querystring receive from URL.Query().
//
// CustomParam args accept bson.M{}, where customParam is anything that
// is not included in the base querystring
//
// [Base QueryString: fields,sort,order,page,offset,id,status]
func ParseQuery(q url.Values, customParam bson.M) map[string]any {
	baseMap := baseQuery(q)

	baseFilter := baseMap["filters"].(bson.M)
	filter := baseFilter
	for k, v := range customParam {
		if v == "array" {
			arr, ok := q[k]
			if !ok {
				continue
			}
			var newArr []any
			var arrK string

			for _, str := range arr {
				id, err := bson.ObjectIDFromHex(str)
				if err != nil {
					newArr = append(newArr, str)
					arrK = k
				} else {
					newArr = append(newArr, id)
					arrK = "_id"
				}
			}
			filter[arrK] = bson.M{"$in": newArr}
		} else {
			if ok := q.Has(k); ok {
				filter[k] = q.Get(k)
			}
		}
	}

	baseMap["filters"] = filter

	return baseMap
}

// MongoFindOpts for repository layer for mongo.Find() arguments.
// Opts can be nil
func MongoFindOpts(filterMap map[string]any, customFilter bson.M) (bson.M, **options.FindOptionsBuilder) {
	var opts *options.FindOptionsBuilder

	filters := bson.M{}
	base := bson.M{}

	if val, ok := filterMap["filters"].(bson.M); ok {
		base = val
	}

	for _, m := range []bson.M{customFilter, base} {
		maps.Copy(filters, m)
	}

	if fields, ok := filterMap["fields"].([]string); ok {
		projection := bson.M{}
		for _, v := range fields {
			projection[v] = 1
		}
		opts = options.Find().SetProjection(projection)
	}
	if page, ok := filterMap["page"].(int64); ok {
		opts = options.Find().SetLimit(page)
	}
	if sort, ok := filterMap["sort"].(map[string]int); ok {
		opts = options.Find().SetSort(sort)
	}

	return filters, &opts
}
