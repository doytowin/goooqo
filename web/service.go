package web

import (
	"context"
	. "github.com/doytowin/goquery/core"
	"regexp"
)

type RestAPI[E any, Q GoQuery] interface {
	Page(ctx context.Context, query Q) (PageList[E], error)
	Get(ctx context.Context, id any) (*E, error)
	CreateMulti(ctx context.Context, entities []E) (int64, error)
	Update(ctx context.Context, e E) (int64, error)
	Patch(ctx context.Context, e E) (int64, error)
	Delete(ctx context.Context, id string) (any, error)
}

type Service[E any, Q GoQuery] struct {
	dataAccess   DataAccess[context.Context, E]
	createQuery  func() Q
	createEntity func() E
	idRgx        *regexp.Regexp
}

func (s *Service[E, Q]) Page(ctx context.Context, query Q) (PageList[E], error) {
	return s.dataAccess.Page(ctx, query)
}

func (s *Service[E, Q]) Get(ctx context.Context, id any) (*E, error) {
	return s.dataAccess.Get(ctx, id)
}

func (s *Service[E, Q]) CreateMulti(ctx context.Context, entities []E) (int64, error) {
	return s.dataAccess.CreateMulti(ctx, entities)
}

func (s *Service[E, Q]) Update(ctx context.Context, entity E) (int64, error) {
	return s.dataAccess.Update(ctx, entity)
}

func (s *Service[E, Q]) Patch(ctx context.Context, entity E) (int64, error) {
	return s.dataAccess.Patch(ctx, entity)
}

func (s *Service[E, Q]) Delete(ctx context.Context, id string) (any, error) {
	return s.dataAccess.Delete(ctx, id)
}

func BuildService[E any, Q GoQuery](
	prefix string,
	dataAccess DataAccess[context.Context, E],
	createEntity func() E,
	createQuery func() Q,
) *Service[E, Q] {
	return &Service[E, Q]{
		dataAccess:   dataAccess,
		createQuery:  createQuery,
		createEntity: createEntity,
		idRgx:        regexp.MustCompile(prefix + `(\d+)$`),
	}
}
