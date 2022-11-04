package idgen_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/utils/idgen"
	"github.com/gavrilaf/wardrobe/pkg/utils/timex"
)

func TestSnowflake(t *testing.T) {

	t.Run("generate sequence ids", func(t *testing.T) {
		g := idgen.NewSnowflake(1)

		id1, err := g.NextID()
		assert.NoError(t, err)
		assert.NotZero(t, id1)

		time.Sleep(time.Millisecond * 2)

		id2, _ := g.NextID()
		assert.Less(t, id1, id2)
	})

	t.Run("generate sequence ids for in one millisecond", func(t *testing.T) {
		timeNow := idgen.SnowflakeTime
		idgen.SnowflakeTime = func() time.Time {
			return timex.Date(2022, 1, 1)
		}

		g := idgen.NewSnowflake(1)

		id1, err := g.NextID()
		assert.NoError(t, err)
		assert.NotZero(t, id1)

		id2, _ := g.NextID()
		assert.Less(t, id1, id2)

		idgen.SnowflakeTime = timeNow
	})

	t.Run("error when clock moving backward", func(t *testing.T) {
		timeNow := idgen.SnowflakeTime
		idgen.SnowflakeTime = func() time.Time {
			return timex.DT(2022, 1, 1, 1, 1, 2, 0)
		}

		g := idgen.NewSnowflake(1)

		_, _ = g.NextID()

		idgen.SnowflakeTime = func() time.Time {
			return timex.DT(2022, 1, 1, 1, 1, 1, 0) // one second less
		}

		_, err := g.NextID()
		assert.Error(t, err)

		idgen.SnowflakeTime = timeNow
	})
}
