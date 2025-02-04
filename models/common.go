// common
package models

import (
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/sports/errors"
	"github.com/nf/geocode"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"strings"
	"time"
)

const (
	AccessRate = 1 << iota // 001
	ThumbRate              // 010
	ReviewRate             // 100

	AccessRateMask = 6 // 110
	ThumbRateMask  = 5 // 101
	ReviewRateMask = 3 // 011

	DefaultPageSize = 20
	TimeFormat      = "2006-01-02 15:04:05"
)

const (
	ScorePhysical = "physical"
	ScoreLiteral  = "literal"
	ScoreMental   = "mental"
	ScoreWealth   = "wealth"
)

const (
	RelNone      = ""
	RelFriend    = "friend"
	RelFollowing = "following"
	RelFollower  = "follower"
	RelBlacklist = "blacklist"
)

var (
	databaseName = "sports"
	accountColl  = "accounts"
	//userColl     = "users"
	articleColl = "articles"
	msgColl     = "messages"
	//reviewColl   = "reviews"
	fileColl   = "files"
	recordColl = "records"
	actionColl = "actions"
	groupColl  = "groups"
	eventColl  = "events"
	ruleColl   = "rules"
	txColl     = "txs"
	configColl = "config"
	//rateColl     = "rates"
)

const (
	Satoshi  = 100000000
	MaxLevel = 40
)

var (
	GuestUserPrefix = "guest:"
	MongoAddr       = "localhost:27017"

	levelScores = make([]int64, MaxLevel)
)

func init() {
	initLevelScores()
}

func InsureIndexes() {
	ensureIndex(accountColl, "_id", "password")
	ensureIndex(accountColl, "-props.score")
	ensureIndex(accountColl, "nickname")
	ensureIndex(accountColl, "-reg_time")
	ensureIndex(accountColl, "-lastlogin")
	ensureIndex2D(accountColl, "loc")

	ensureIndex(articleColl, "author")
	ensureIndex(articleColl, "-pub_time")
	ensureIndex(articleColl, "-pub_time", "-_id")

	ensureIndex(eventColl, "-time")

	ensureIndex(fileColl, "fid")
	ensureIndex(fileColl, "-uploadDate")

	ensureIndex(msgColl, "from")
	ensureIndex(msgColl, "to")
	ensureIndex(msgColl, "from", "to")
	ensureIndex(msgColl, "-time")

	ensureIndex(recordColl, "uid")
	ensureIndex(recordColl, "-sport.time")
	ensureIndex(recordColl, "-sport.distance")
	ensureIndex(recordColl, "-pub_time")
	ensureIndex(recordColl, "-pub_time", "-_id")

	ensureIndex(actionColl, "userid")
	ensureIndex(actionColl, "date")
	ensureIndex(actionColl, "userid", "date")

	ensureIndex(txColl, "-time")
	ensureIndex(txColl, "uid")
}

func scoreOfUpgrade(n int) int64 {
	difficult := func(n int) int {
		if n < 10 {
			return 0
		} else if n < 20 {
			return 1
		} else if n < 30 {
			return 3
		} else if n < 35 {
			return 6
		} else {
			return 5 * (n - 33)
		}
	}

	factor := func(n int) float64 {
		if n <= 10 {
			return 1
		} else if n < 30 {
			return (1.0 - float64(n-10)/100)
		} else {
			return 0.82
		}
	}

	s := int64(float64(2*n+difficult(n)) * float64(40+3*n) * factor(n))
	return s - s%10
}

func initLevelScores() {
	var total int64
	for i := 1; i < len(levelScores); i++ {
		total += scoreOfUpgrade(i)
		levelScores[i] = total
		//fmt.Println(i, total)
	}
}

func GetLevelScore() []int64 {
	return levelScores
}

func Score2Level(score int64) int64 {
	for i := 1; i < len(levelScores); i++ {
		if score < levelScores[i] {
			return int64(i)
		}
	}

	return int64(MaxLevel)
}

type Paging struct {
	First string `form:"page_frist_id" json:"page_frist_id"`
	Last  string `form:"page_last_id" json:"page_last_id"`
	Count int    `form:"page_item_count" json:"-"`
}

type Address struct {
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	Area     string `json:"area,omitempty"`
	Desc     string `bson:"location_desc" json:"location_desc"`
}

func (addr Address) String() string {
	return addr.Country + addr.Province + addr.City + addr.Area + addr.Desc
}

type Location struct {
	Lat float64 `bson:"latitude" json:"latitude"`
	Lng float64 `bson:"longitude" json:"longitude"`
}

func Addr2Loc(addr Address) Location {
	return Location{}

	if len(addr.String()) == 0 {
		return Location{}
	}
	req := &geocode.Request{
		Region:   "cn",
		Provider: geocode.GOOGLE,
		Address:  addr.String(),
	}
	resp, err := req.Lookup(nil)
	if err != nil || resp.Status != "OK" || len(resp.Results) == 0 {
		return Location{}
	}

	return Location{
		Lat: resp.Results[0].Geometry.Location.Lat,
		Lng: resp.Results[0].Geometry.Location.Lng,
	}
}

type PagingFunc func(c *mgo.Collection, first, last string, args ...interface{}) (query bson.M, err error)

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MongoAddr)
		//log.Println("mongo addr:", MongoAddr)
		if err != nil {
			log.Println(err) // no, not really
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}

func withCollection(collection string, safe *mgo.Safe, s func(*mgo.Collection) error) error {
	session, err := getSession()
	if session == nil {
		return err
	}
	defer session.Close()

	session.SetSafe(safe)
	c := session.DB(databaseName).C(collection)
	return s(c)
}

func exists(collection string, query interface{}) (bool, error) {
	b := false
	q := func(c *mgo.Collection) error {
		n, err := c.Find(query).Count()
		b = n > 0
		return err
	}

	err := withCollection(collection, nil, q)
	return b, err
}

func query(collection string, qry, selector interface{},
	sortFields []string, limit int) (*mgo.Query, error) {
	session, err := getSession()
	if session == nil {
		return nil, err
	}
	defer session.Close()

	c := session.DB(databaseName).C(collection)

	return c.Find(qry).Select(selector).Sort(sortFields...).Limit(limit), nil
}

// search with paging
func psearch(collection string, query, selector interface{}, sortFields []string,
	total *int, result interface{}, pagingFunc PagingFunc, paging *Paging, args ...interface{}) (err error) {

	q := func(c *mgo.Collection) error {
		var pquery bson.M
		if pagingFunc != nil {
			if paging == nil {
				paging = &Paging{}
			}
			pquery, err = pagingFunc(c, paging.First, paging.Last, args...)
			if err != nil {
				return err
			}
		}

		qy := c.Find(query)

		if total != nil {
			if *total, err = qy.Count(); err != nil {
				return err
			}
		}
		if result == nil {
			return nil
		}

		if paging.Count == 0 {
			paging.Count = DefaultPageSize
		}

		if pquery == nil {
			pquery = bson.M{}
		}
		if m, ok := query.(bson.M); ok {
			for k, v := range m {
				pquery[k] = v
			}
		}

		return c.Find(pquery).Select(selector).Sort(sortFields...).Limit(paging.Count).All(result)
	}
	return withCollection(collection, nil, q)
}

func search(collection string, query interface{}, selector interface{},
	skip, limit int, sortFields []string, total *int, result interface{}) error {

	q := func(c *mgo.Collection) error {
		qy := c.Find(query)
		var err error

		if selector != nil {
			qy = qy.Select(selector)
		}

		if total != nil {
			if *total, err = qy.Count(); err != nil {
				return err
			}
		}

		if result == nil {
			return err
		}

		if limit > 0 {
			qy = qy.Limit(limit)
		}
		if skip > 0 {
			qy = qy.Skip(skip)
		}
		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.All(result)
	}

	if err := withCollection(collection, nil, q); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func count(collection string, query interface{}) (count int, err error) {
	q := func(c *mgo.Collection) (err error) {
		count, err = c.Find(query).Count()
		return
	}

	if e := withCollection(collection, nil, q); e != nil {
		err = errors.NewError(errors.DbError)
	}
	return
}

func findOne(collection string, query interface{}, sortFields []string, result interface{}) error {
	q := func(c *mgo.Collection) error {
		var err error
		qy := c.Find(query)

		if result == nil {
			return err
		}

		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.One(result)
	}

	if err := withCollection(collection, nil, q); err != nil {
		return errors.NewError(errors.DbError)
	}
	return nil
}

func findIds(c string, ids []interface{}, result interface{}) error {
	return search(c, bson.M{"_id": bson.M{"$in": ids}}, nil, 0, 0, nil, nil, result)
}

func updateId(collection string, id interface{}, change interface{}, safe bool) error {
	update := func(c *mgo.Collection) error {
		return c.UpdateId(id, change)
	}

	if safe {
		return withCollection(collection, &mgo.Safe{}, update)
	}
	return withCollection(collection, nil, update)
}

func update(collection string, selector, change interface{}, safe bool) error {
	update := func(c *mgo.Collection) error {
		return c.Update(selector, change)
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, update)
	}
	return withCollection(collection, nil, update)
}

func updateAll(collection string, selector, change interface{}, safe bool) (info *mgo.ChangeInfo, err error) {
	update := func(c *mgo.Collection) error {
		info, err = c.UpdateAll(selector, change)
		return err
	}
	if safe {
		withCollection(collection, &mgo.Safe{}, update)
	} else {
		withCollection(collection, nil, update)
	}
	if err != nil {
		return info, errors.NewError(errors.DbError)
	}

	return
}

func upsert(collection string, selector, change interface{}, safe bool) (*mgo.ChangeInfo, error) {
	var chinfo *mgo.ChangeInfo

	upsert := func(c *mgo.Collection) (err error) {
		chinfo, err = c.Upsert(selector, change)
		//log.Println(chinfo, err)
		return err
	}
	if safe {
		return chinfo, withCollection(collection, &mgo.Safe{}, upsert)
	}
	return chinfo, withCollection(collection, nil, upsert)
}

func save(collection string, o interface{}, safe bool) error {
	insert := func(c *mgo.Collection) error {
		return c.Insert(o)
	}

	if safe {
		return withCollection(collection, &mgo.Safe{}, insert)
	}
	return withCollection(collection, nil, insert)
}

func remove(collection string, selector interface{}, safe bool) error {
	rm := func(c *mgo.Collection) error {
		return c.Remove(selector)
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, rm)
	}
	return withCollection(collection, nil, rm)
}

func removeId(collection string, id interface{}, safe bool) error {
	rm := func(c *mgo.Collection) error {
		return c.RemoveId(id)
	}
	if safe {
		return withCollection(collection, &mgo.Safe{}, rm)
	}
	return withCollection(collection, nil, rm)
}

func removeAll(collection string, selector interface{}, safe bool) (info *mgo.ChangeInfo, err error) {
	r := func(c *mgo.Collection) error {
		info, err = c.RemoveAll(selector)
		return err
	}
	if safe {
		withCollection(collection, &mgo.Safe{}, r)
	} else {
		withCollection(collection, nil, r)
	}
	if err != nil {
		return info, errors.NewError(errors.DbError)
	}

	return
}

func apply(collection string, selector interface{}, change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	apply := func(c *mgo.Collection) (err error) {
		info, err = c.Find(selector).Apply(change, result)
		return err
	}

	err = withCollection(collection, nil, apply)
	return
}

func ensureIndex(collection string, keys ...string) error {
	ensure := func(c *mgo.Collection) error {
		return c.EnsureIndexKey(keys...)
	}

	return withCollection(collection, nil, ensure)
}

func ensureIndex2D(collection string, key string) error {
	ensure := func(c *mgo.Collection) error {
		return c.EnsureIndex(mgo.Index{
			Key: []string{"$2d:" + key},
		})
	}
	return withCollection(collection, nil, ensure)
}

func DateString(t time.Time) string {
	return t.Format("2006-01-02")
}

func Struct2Map(i interface{}) bson.M {
	v, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	fmt.Println("(struct2map):", string(v))
	var m bson.M
	json.Unmarshal(v, &m)

	return m
}

func AgeToTimeRange(age int) (start time.Time, end time.Time) {
	now := time.Now()
	start = time.Date(now.Year()-age, time.January, 1, 0, 0, 0, 0, now.Location())
	end = time.Date(now.Year()-age, time.December, 31, 23, 59, 59, 999999, now.Location())

	return
}

func charFilter(s string) string {
	chars := ".*$^?"
	if !strings.ContainsAny(s, chars) {
		return s
	}

	r := strings.NewReplacer(".", "\\.",
		"*", "\\*",
		"$", "\\$",
		"^", "\\^")

	return r.Replace(s)
}
