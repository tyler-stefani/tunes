package cache

import "testing"

func TestCacheIntegration(t *testing.T) {
	c := NewCache()

	c.Put("testKey", "testData")

	v := c.Get("testKey")

	if v != "testData" {
		t.Fatalf("got incorrect value")
	}
}
