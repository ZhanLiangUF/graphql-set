package pg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	_ "github.com/lib/pq"
)

type Repository interface {
	CreateSet(ctx context.Context, setDatas []int64) (*Set, map[string][]int64, error)
	ListSetsWithIntersectingSets(ctx context.Context) (map[string][]int64, map[string][]string, error)
}

type repoSvc struct {
	*Queries
	db *sql.DB
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

func (r *repoSvc) CreateSet(ctx context.Context, setDatas []int64) (*Set, map[string][]int64, error) {
	// get hash of set
	var buffer bytes.Buffer
	sort.Slice(setDatas, func(i, j int) bool {
		return setDatas[i] < setDatas[j]
	})
	for i, _ := range setDatas {
		buffer.WriteString(strconv.FormatInt(setDatas[i], 36))
		buffer.WriteString(" ")
	}
	checksum := sha256.Sum256(buffer.Bytes())
	// insert into sets and setsdatas table
	set := new(Set)
	m := make(map[string][]int64)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.CreateSet(ctx, checksum[:])
		if err != nil {
			return err
		}
		if len(setDatas) == 0 {
			if err := q.CreateSetData(ctx, CreateSetDataParams{
				Data:   sql.NullInt64{Valid: false},
				SetUid: res.SetUid,
			}); err != nil {
				return err
			}
		} else {
			for _, setDatum := range setDatas {
				if err := q.CreateSetData(ctx, CreateSetDataParams{
					Data: sql.NullInt64{
						Int64: setDatum,
						Valid: true,
					},
					SetUid: res.SetUid,
				}); err != nil {
					return err
				}
			}
		}
		// create map of set_uid to their data and then go through each unique set
		// and insert into intersecting sets if one of the data matches using binary search
		var intersectingSets [][]byte

		allSets, err := q.ListSetsDatas(ctx)
		if err != nil {
			return err
		}

		for _, setData := range allSets {

			id := hex.EncodeToString(setData.SetUid)
			if val, ok := m[id]; ok {
				m[id] = append(val, setData.Data.Int64)
			} else {
				if !setData.Data.Valid {
					m[id] = []int64{}
				} else {
					m[id] = []int64{setData.Data.Int64}
				}
			}
			if !contains(intersectingSets, setData.SetUid) {
				i := sort.Search(len(setDatas), func(i int) bool { return setDatas[i] == setData.Data.Int64 })
				if setDatas[i] == setData.Data.Int64 {
					intersectingSets = append(intersectingSets, setData.SetUid)
				}
			}
		}
		for _, setuid := range intersectingSets {
			err := q.SetIntersectingSet(ctx, SetIntersectingSetParams{SetUid: res.SetUid, IntersectingsetUid: setuid})
			if err != nil {
				return err
			}
		}
		for k, _ := range m {
			decodeK, _ := hex.DecodeString(k)
			if !contains(intersectingSets, decodeK) {
				delete(m, k)
			}
		}
		set = &res
		return nil
	})
	return set, m, err
}

func (r *repoSvc) ListSetsWithIntersectingSets(ctx context.Context) (map[string][]int64, map[string][]string, error) {
	smap := make(map[string][]int64)
	ismap := make(map[string][]string)
	err := r.withTx(ctx, func(q *Queries) error {
		allSets, err := q.ListSetsDatas(ctx)
		if err != nil {
			return err
		}

		for _, setData := range allSets {
			id := hex.EncodeToString(setData.SetUid)
			if val, ok := smap[id]; ok {
				smap[id] = append(val, setData.Data.Int64)
			} else {
				if !setData.Data.Valid {
					smap[id] = []int64{}
				} else {
					smap[id] = []int64{setData.Data.Int64}
				}
			}
		}

		intersectingSets, err := q.ListIntersectingSets(ctx)
		if err != nil {
			return err
		}
		for _, is := range intersectingSets {
			id := hex.EncodeToString(is.SetUid)
			isid := hex.EncodeToString(is.IntersectingsetUid)
			if val, ok := ismap[id]; ok {
				ismap[id] = append(val, isid)
			} else {
				ismap[id] = []string{isid}
			}
		}
		return nil
	})
	return smap, ismap, err
}

func contains(s [][]byte, t []byte) bool {
	for _, b := range s {
		if bytes.Equal(b, t) {
			return true
		}
	}
	return false
}

func NewRepository(db *sql.DB) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}
