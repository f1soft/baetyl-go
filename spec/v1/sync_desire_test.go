package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRDData(t *testing.T) {
	{
		// --- app
		desireddata := &ResourceValue{}
		desireddata.Name = "app"
		desireddata.Version = "123"
		desireddata.Kind = KindApplication
		desireddata.Value.Value = &Application{Name: "c"}
		expected := "{\"name\":\"c\",\"createTime\":\"0001-01-01T00:00:00Z\"}"

		appdata, err := json.Marshal(desireddata)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(desireddata.Value.Data))
		fmt.Printf(string(appdata))

		desireddata2 := &ResourceValue{}
		err = json.Unmarshal(appdata, desireddata2)
		assert.NoError(t, err)
		assert.Nil(t, desireddata2.Value.Value)
		assert.Equal(t, expected, string(desireddata.Value.Data))

		// success
		app := desireddata2.App()
		assert.Equal(t, desireddata.Value.Value, app)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)

		desireddata.Kind = KindApp
		app = desireddata2.App()
		assert.Equal(t, desireddata.Value.Value, app)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)

		// failure
		cfg := desireddata2.Config()
		assert.Nil(t, cfg)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)

		// failure
		scr := desireddata2.Secret()
		assert.Nil(t, scr)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)
	}
	{
		// --- config
		desireddata := &ResourceValue{}
		desireddata.Name = "cfg"
		desireddata.Version = "123"
		desireddata.Kind = KindConfiguration
		desireddata.Value.Value = &Configuration{Name: "c"}
		expected := "{\"name\":\"c\",\"createTime\":\"0001-01-01T00:00:00Z\",\"updateTime\":\"0001-01-01T00:00:00Z\"}"

		appdata, err := json.Marshal(desireddata)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(desireddata.Value.Data))
		fmt.Printf(string(appdata))

		desireddata2 := &ResourceValue{}
		err = json.Unmarshal(appdata, desireddata2)
		assert.NoError(t, err)
		assert.Nil(t, desireddata2.Value.Value)
		assert.Equal(t, expected, string(desireddata.Value.Data))

		// failure
		app := desireddata2.App()
		assert.Nil(t, app)
		assert.Nil(t, desireddata2.Value.Value)

		// sucees
		cfg := desireddata2.Config()
		assert.Equal(t, desireddata.Value.Value, cfg)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)

		desireddata.Kind = KindConfig
		cfg = desireddata2.Config()
		assert.Equal(t, desireddata.Value.Value, cfg)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)

		// failure
		scr := desireddata2.Secret()
		assert.Nil(t, scr)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)
	}
	{
		// --- secret
		desireddata := &ResourceValue{}
		desireddata.Name = "scr"
		desireddata.Version = "123"
		desireddata.Kind = KindSecret
		desireddata.Value.Value = &Secret{Name: "c"}
		expected := "{\"name\":\"c\",\"createTime\":\"0001-01-01T00:00:00Z\",\"updateTime\":\"0001-01-01T00:00:00Z\"}"

		appdata, err := json.Marshal(desireddata)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(desireddata.Value.Data))
		fmt.Printf(string(appdata))

		desireddata2 := &ResourceValue{}
		err = json.Unmarshal(appdata, desireddata2)
		assert.NoError(t, err)
		assert.Nil(t, desireddata2.Value.Value)
		assert.Equal(t, expected, string(desireddata.Value.Data))

		// failure
		app := desireddata2.App()
		assert.Nil(t, app)
		assert.Nil(t, desireddata2.Value.Value)

		// failure
		cfg := desireddata2.Config()
		assert.Nil(t, cfg)
		assert.Nil(t, desireddata2.Value.Value)

		// failure
		scr := desireddata2.Secret()
		assert.Equal(t, desireddata.Value.Value, scr)
		assert.Equal(t, desireddata.Value.Value, desireddata2.Value.Value)
	}
}
