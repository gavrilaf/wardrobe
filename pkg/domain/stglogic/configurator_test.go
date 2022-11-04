package stglogic_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/domain/stglogic"
	"github.com/gavrilaf/wardrobe/pkg/fs/fsmocks"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
	"github.com/gavrilaf/wardrobe/pkg/utils/idgen/idgenmocks"
	"github.com/gavrilaf/wardrobe/pkg/utils/timex"
)

func TestGenerateFileName(t *testing.T) {

	obj := dbtypes.InfoObject{
		Published: timex.DT(2022, 10, 23, 12, 5, 0, 0),
	}

	t.Run("should keep extension", func(t *testing.T) {
		subj, mm := subjWithMocks()

		mm.snf.On("NextID").Return(int64(123), nil)

		t.Run("with extension", func(t *testing.T) {
			fn, err := subj.GenerateFileName(obj, dbtypes.File{Name: "test.jpg"})
			assert.NoError(t, err)
			assert.Equal(t, "2022-10-23-12-05-test-123.jpg", fn)
		})

		t.Run("only dot", func(t *testing.T) {
			fn, err := subj.GenerateFileName(obj, dbtypes.File{Name: "test."})
			assert.NoError(t, err)
			assert.Equal(t, "2022-10-23-12-05-test-123.", fn)
		})

		t.Run("no extension", func(t *testing.T) {
			fn, err := subj.GenerateFileName(obj, dbtypes.File{Name: "test"})
			assert.NoError(t, err)
			assert.Equal(t, "2022-10-23-12-05-test-123", fn)
		})
	})

	t.Run("failed when snowflake failed", func(t *testing.T) {
		subj, mm := subjWithMocks()

		mm.snf.On("NextID").Return(int64(0), fmt.Errorf(""))

		_, err := subj.GenerateFileName(obj, dbtypes.File{Name: "test.jpg"})
		assert.Error(t, err)
	})
}

type mmocks struct {
	fs  *fsmocks.Storage
	snf *idgenmocks.Snowflake
}

func subjWithMocks() (stglogic.Configurator, *mmocks) {
	mm := &mmocks{
		fs:  &fsmocks.Storage{},
		snf: &idgenmocks.Snowflake{},
	}

	return stglogic.NewConfigurator(mm.fs, mm.snf), mm
}
