package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_memoryClient_Get(t *testing.T) {
	client := NewMemoryClient()

	item, err := client.Get("key")
	assert.Nil(t, item)
	assert.NoError(t, err)
}

func Test_memoryClient_Set(t *testing.T) {
	client := NewMemoryClient()

	err := client.Set("key", &Item{Data: "foo"})
	assert.NoError(t, err)

	item, err := client.Get("key")
	assert.Equal(t, item.Data.(string), "foo")
	assert.NoError(t, err)
}

func Test_memoryClient_Add(t *testing.T) {
	client := NewMemoryClient()

	err := client.Add("key", &Item{Data: "foo"})
	assert.NoError(t, err)

	item, err := client.Get("key")
	assert.Equal(t, item.Data.(string), "foo")
	assert.NoError(t, err)

	err = client.Add("key", &Item{Data: "foo"})
	assert.EqualError(t, err, "not stored")
}

func Test_memoryClient_Delete(t *testing.T) {
    client := NewMemoryClient()

    err := client.Set("key", &Item{Data: "foo"})
    assert.NoError(t, err)

    item, err := client.Get("key")
    assert.Equal(t, item.Data.(string), "foo")
    assert.NoError(t, err)

    err = client.Delete("key")
    assert.NoError(t, err)

    item, err = client.Get("key")
    assert.Nil(t, item)
    assert.NoError(t, err)

    err = client.Delete("key")
    assert.NoError(t, err)
}
