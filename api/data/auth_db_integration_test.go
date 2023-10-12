package data

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAuthDBIntegration(t *testing.T) {
	adb, err := NewAuthDB(os.Getenv("AUTH_DATABASE_URL"))
	if err != nil {
		t.Skip("skipping integration test")
	}

	expiration := time.Now().Add(time.Hour * 24).Truncate(time.Second)
	expected := RefreshToken{"test", expiration}

	if ok, err := adb.WriteRefreshToken("test", expiration); !ok {
		t.Error(err)
	}

	if actual, err := adb.FindRefreshToken("test"); err != nil {
		t.Error(err)
	} else {
		if reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected %+v but got %+v", expected, actual)
		}
	}
}
