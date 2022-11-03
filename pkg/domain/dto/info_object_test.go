package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
	"github.com/gavrilaf/wardrobe/pkg/utils"
)

func TestInfoObject(t *testing.T) {
	dbObj := dbtypes.InfoObject{
		Name:      "name",
		Author:    "author",
		Source:    "source",
		Published: utils.DT(2002, time.April, 23, 12, 45, 0, 0),
	}

	dtoObj := dto.InfoObjectFromDBType(dbObj)

	dbObj2, err := dtoObj.ToDBType()
	assert.NoError(t, err)

	assert.Equal(t, dbObj, dbObj2)

}
