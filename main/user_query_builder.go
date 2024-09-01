package main

import "github.com/doytowin/goooqo/rdb"
import "strings"

func (q UserQuery) BuildConditions() ([]string, []any) {
	conditions := make([]string, 0, 4)
	args := make([]any, 0, 4)
	if q.IdGt != nil {
		conditions = append(conditions, "id > ?")
		args = append(args, q.IdGt)
	}
	if q.IdIn != nil {
		phs := make([]string, 0, len(*q.IdIn))
		for _, arg := range *q.IdIn {
			args = append(args, arg)
			phs = append(phs, "?")
		}
		conditions = append(conditions, "id IN ("+strings.Join(phs, ", ")+")")
	}
	if q.IdNotIn != nil {
		phs := make([]string, 0, len(*q.IdNotIn))
		for _, arg := range *q.IdNotIn {
			args = append(args, arg)
			phs = append(phs, "?")
		}
		conditions = append(conditions, "id NOT IN ("+strings.Join(phs, ", ")+")")
	}
	if q.Cond != nil {
		conditions = append(conditions, "(Score = ? OR Memo = ?)")
		args = append(args, q.Cond)
		args = append(args, q.Cond)
	}
	if q.ScoreLt != nil {
		conditions = append(conditions, "score < ?")
		args = append(args, q.ScoreLt)
	}
	if q.ScoreLt1 != nil {
		whereClause, args1 := rdb.BuildWhereClause(q.ScoreLt1)
		condition := "score < (SELECT avg(score) FROM User" + whereClause + ")"
		conditions = append(conditions, condition)
		args = append(args, args1...)
	}
	if q.MemoNull != nil {
		if *q.MemoNull {
			conditions = append(conditions, "memo IS NULL")
		} else {
			conditions = append(conditions, "memo IS NOT NULL")
		}
	}
	if q.MemoLike != nil {
		conditions = append(conditions, "memo LIKE ?")
		args = append(args, q.MemoLike)
	}
	if q.Deleted != nil {
		conditions = append(conditions, "deleted = ?")
		args = append(args, q.Deleted)
	}
	return conditions, args
}
