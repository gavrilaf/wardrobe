package idgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/utils/idgen"
)

func TestNodeID(t *testing.T) {
	nodeID, err := idgen.NodeID()
	assert.NoError(t, err)
	assert.Greater(t, nodeID, int64(0))
}

func TestNormalizeNodeID(t *testing.T) {
	assert.Equal(t, int64(123), idgen.NormalizeNodeID(123))
	assert.Equal(t, int64(idgen.MaxNodeID), idgen.NormalizeNodeID(idgen.MaxNodeID))

	assert.Equal(t, int64(561), idgen.NormalizeNodeID(uint64(idgen.MaxNodeID)+100))
}
