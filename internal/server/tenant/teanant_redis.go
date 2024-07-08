package tenant

import (
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/go-redis/redis"
)

type TenantCache struct {
	R *redis.Client
	W *redis.Client
}

func (c *TenantCache) ZKey() string {
	return "tenant"
}

func (c *TenantCache) Key(uuid int) string {
	return fmt.Sprintf("%s-%d", c.ZKey(), uuid)
}

func (c *TenantCache) ZAdd(ms *[]server.TenantMeta) error {
	zs := make([]redis.Z, len(*ms))
	for i, m := range *ms {
		zs[i] = redis.Z{Score: float64(m.Uuid), Member: m.Name}
	}
	_, err := c.W.ZAdd(c.ZKey(), zs...).Result()

	return err
}

func (c *TenantCache) ZRank(member string) (int, error) {
	result, err := c.R.ZRank(c.ZKey(), member).Result()
	return int(result), err
}

func (c *TenantCache) ZScore(member string) (int, error) {
	result, err := c.R.ZScore(c.ZKey(), member).Result()
	return int(result), err
}

func (c *TenantCache) ZRem(member string) error {
	_, err := c.W.ZRem(c.ZKey(), member).Result()
	return err
}

func (c *TenantCache) Set(m *server.TenantMeta) error {
	return c.W.Set(c.Key(m.Uuid), m, 0).Err()
}

func (c *TenantCache) Get(uuid int) (t server.TenantMeta, err error) {
	err = c.R.Get(c.Key(uuid)).Scan(&t)
	return t, err
}

func (c *TenantCache) Del(uuid int) error {
	_, err := c.W.Del(c.Key(uuid)).Result()
	return err
}

func (c *TenantCache) ZAddSet(ms *[]server.TenantMeta) (err error) {
	err = c.ZAdd(ms)
	if err != nil {
		return err
	}
	for _, m := range *ms {
		err = c.Set(&m)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *TenantCache) ZRemDel(member string, uuid int) (err error) {
	err = c.ZRem(member)
	if err != nil {
		return err
	}

	err = c.Del(uuid)
	return err
}

func (c *TenantCache) ZScoreGet(member string) (t server.TenantMeta, err error) {
	tid, err := c.ZScore(member)
	if err != nil {
		return t, err
	}

	t, err = c.Get(tid)
	return t, err
}
