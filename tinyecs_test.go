package tinyecs_test

import (
	"github.com/kaiaverkvist/tinyecs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type playerData struct {
	name   string
	health float32
}

type floater struct {
	f float64
}

type velocity struct {
	v float64
}

type testEntity struct {
	tinyecs.Entity

	name string
}

func Test_BasicInitialization(t *testing.T) {
	e := tinyecs.NewEngine()

	entity := testEntity{}
	e.AddComponents(
		entity,

		floater{f: 1.0},
		velocity{v: 5.0},
		floater{1},
	)

	e.AddEntity(&entity)

	assert.Len(t, e.GetComponents(), 3)
}

func Test_RemoveComponentById(t *testing.T) {
	e := tinyecs.NewEngine()

	entity := testEntity{}
	e.AddComponents(
		entity,

		floater{f: 1.0},
		velocity{v: 5.0},
		floater{1},
	)
	e.AddEntity(&entity)

	assert.Len(t, e.GetComponents(), 3)

	e.DeleteComponent(e.GetComponents()[1])

	assert.Len(t, e.GetComponents(), 2)
}

func TestEngine_AddEntity(t *testing.T) {
	e := tinyecs.NewEngine()
	e.GetEntities()

	assert.Len(t, e.GetEntities(), 0)

	entity := testEntity{}
	e.AddComponents(
		entity,

		floater{f: 1.0},
		velocity{v: 5.0},
		floater{1},
	)

	e.AddEntity(&entity)

	assert.Len(t, e.GetEntities(), 1)
}

func TestEngine_RemoveEntity(t *testing.T) {
	e := tinyecs.NewEngine()
	e.GetEntities()

	assert.Len(t, e.GetEntities(), 0)

	entity := testEntity{}
	e.AddComponents(
		entity,

		floater{f: 1.0},
		velocity{v: 5.0},
		floater{1},
	)

	e.AddEntity(&entity)

	assert.Len(t, e.GetEntities(), 1)

	e.RemoveEntity(&entity)

	assert.Len(t, e.GetEntities(), 0)
}

func Test_EachAndUpdateInstances(t *testing.T) {
	e := tinyecs.NewEngine()

	entity := testEntity{}
	e.AddComponents(
		entity,

		floater{f: 1.0},
		velocity{v: 5.0},
		floater{1},
	)

	e.AddEntity(&entity)

	tinyecs.Each[floater](&e, func(id uint64, obj floater) {
		assert.Equal(t, 1.0, obj.f)
		obj.f += 0.35
		tinyecs.Set(&e, id, obj)
	})

	tinyecs.Each[floater](&e, func(id uint64, obj floater) {
		assert.Equal(t, 1.35, obj.f)
	})
}

func Test_EntityWithComponents(t *testing.T) {
	e := tinyecs.NewEngine()

	entity := testEntity{}
	e.AddComponents(
		entity,

		velocity{},
		playerData{name: "test", health: 100.0},
	)
	e.AddEntity(&entity)

	c := tinyecs.Each[velocity](&e, func(id uint64, obj velocity) {})
	assert.Equal(t, uint64(1), c)

	tinyecs.Each[playerData](&e, func(id uint64, obj playerData) {
		assert.Equal(t, "test", obj.name)
	})
}

func Test_EachEntity(t *testing.T) {
	e := tinyecs.NewEngine()

	entity := testEntity{name: "a1"}
	e.AddComponents(
		entity,

		velocity{},
		playerData{name: "test", health: 100.0},
	)
	e.AddEntity(&entity)

	c := tinyecs.EachEntity[testEntity, velocity](&e, func(entity testEntity, component velocity) {
		assert.Equal(t, entity.name, "a1")
	})
	assert.Equal(t, uint64(1), c)

}
