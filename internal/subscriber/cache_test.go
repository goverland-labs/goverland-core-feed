package subscriber

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUnitNewCache(t *testing.T) {
	c := NewCache()
	require.NotNil(t, c)
	require.NotNil(t, c.data)
}

func TestUnitUpsertAndGetItems(t *testing.T) {
	c := NewCache()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	c.UpsertItem(id1, &Subscriber{ID: id1, WebhookURL: "url-1"})
	c.UpsertItem(id2, &Subscriber{ID: id2, WebhookURL: "url-2"})

	t.Run("contains correct values", func(t *testing.T) {
		value, ok := c.GetItem(id1)
		require.True(t, ok)
		require.Equal(t, &Subscriber{ID: id1, WebhookURL: "url-1"}, value)
	})

	t.Run("do not collect unexpected values", func(t *testing.T) {
		value, ok := c.GetItem(uuid.New())
		require.False(t, ok)
		require.Empty(t, value)
	})

	t.Run("upsert change value", func(t *testing.T) {
		c.UpsertItem(id1, &Subscriber{ID: id3, WebhookURL: "url-3"})
		value, ok := c.GetItem(id1)
		require.True(t, ok)
		require.Equal(t, &Subscriber{ID: id3, WebhookURL: "url-3"}, value)
	})
}

func TestUnitRemoveItem(t *testing.T) {
	c := NewCache()

	id := uuid.New()

	c.UpsertItem(id, &Subscriber{})
	_, ok := c.GetItem(id)
	require.True(t, ok)

	c.RemoveItem(id)

	_, ok = c.GetItem(id)
	require.False(t, ok)
}
