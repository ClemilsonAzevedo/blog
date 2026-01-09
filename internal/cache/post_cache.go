package cache

import (
	"fmt"
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"go.bryk.io/pkg/ulid"
)

const (
	postByIDPrefix   = "post:id:"
	postBySlugPrefix = "post:slug:"
	allPostsKey      = "posts:all"
	paginatedPrefix  = "posts:paginated:"
)

type PostCache struct {
	singlePostCache *Cache[*entities.Post]
	listCache       *Cache[[]*entities.Post]
	paginatedCache  *Cache[PaginatedResult]
}

type PaginatedResult struct {
	Posts []entities.Post
	Total int64
}

func NewPostCache(ttl time.Duration) *PostCache {
	return &PostCache{
		singlePostCache: NewCache[*entities.Post](ttl),
		listCache:       NewCache[[]*entities.Post](ttl),
		paginatedCache:  NewCache[PaginatedResult](ttl),
	}
}

func (pc *PostCache) GetByID(id ulid.ULID) (*entities.Post, bool) {
	key := postByIDPrefix + id.String()
	return pc.singlePostCache.Get(key)
}

func (pc *PostCache) SetByID(id ulid.ULID, post *entities.Post) {
	key := postByIDPrefix + id.String()
	pc.singlePostCache.Set(key, post)
}

func (pc *PostCache) GetBySlug(slug string) (*entities.Post, bool) {
	key := postBySlugPrefix + slug
	return pc.singlePostCache.Get(key)
}

func (pc *PostCache) SetBySlug(slug string, post *entities.Post) {
	key := postBySlugPrefix + slug
	pc.singlePostCache.Set(key, post)
}

func (pc *PostCache) GetAll() ([]*entities.Post, bool) {
	return pc.listCache.Get(allPostsKey)
}

func (pc *PostCache) SetAll(posts []*entities.Post) {
	pc.listCache.Set(allPostsKey, posts)
}

func (pc *PostCache) GetPaginated(page, limit int) (PaginatedResult, bool) {
	key := fmt.Sprintf("%s%d:%d", paginatedPrefix, page, limit)
	return pc.paginatedCache.Get(key)
}

func (pc *PostCache) SetPaginated(page, limit int, posts []entities.Post, total int64) {
	key := fmt.Sprintf("%s%d:%d", paginatedPrefix, page, limit)
	pc.paginatedCache.Set(key, PaginatedResult{
		Posts: posts,
		Total: total,
	})
}

func (pc *PostCache) InvalidatePost(id ulid.ULID, slug string) {
	pc.singlePostCache.Delete(postByIDPrefix + id.String())
	if slug != "" {
		pc.singlePostCache.Delete(postBySlugPrefix + slug)
	}
}

func (pc *PostCache) InvalidateLists() {
	pc.listCache.Clear()
	pc.paginatedCache.Clear()
}

func (pc *PostCache) InvalidateAll() {
	pc.singlePostCache.Clear()
	pc.listCache.Clear()
	pc.paginatedCache.Clear()
}
