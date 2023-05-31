package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNewCache(t *testing.T) {
	c := NewCache()
	require.NotNil(t, c)
	require.NotNil(t, c.data)
}

func TestUnitUpsertAndGetItems(t *testing.T) {
	c := NewCache()
	c.UpsertItem("key-1", &Subscriber{ID: "1", WebhookURL: "url-1"})
	c.UpsertItem("key-2", &Subscriber{ID: "2", WebhookURL: "url-2"})

	t.Run("contains correct values", func(t *testing.T) {
		value, ok := c.GetItem("key-1")
		require.True(t, ok)
		require.Equal(t, &Subscriber{ID: "1", WebhookURL: "url-1"}, value)
	})

	t.Run("do not collect unexpected values", func(t *testing.T) {
		value, ok := c.GetItem("unknown")
		require.False(t, ok)
		require.Empty(t, value)
	})

	t.Run("upsert change value", func(t *testing.T) {
		c.UpsertItem("key-1", &Subscriber{ID: "3", WebhookURL: "url-3"})
		value, ok := c.GetItem("key-1")
		require.True(t, ok)
		require.Equal(t, &Subscriber{ID: "3", WebhookURL: "url-3"}, value)
	})
}

func TestUnitRemoveItem(t *testing.T) {
	c := NewCache()

	c.UpsertItem("key-1", &Subscriber{})
	_, ok := c.GetItem("key-1")
	require.True(t, ok)

	c.RemoveItem("key-1")

	_, ok = c.GetItem("key-1")
	require.False(t, ok)
}
