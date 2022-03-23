package tinyecs

import (
	"reflect"
	"sync"
)

// entityComponentLink is used to store a relationship between an entity and a component.
type entityComponentLink struct {
	entity    any
	component *any
}

// Engine represents the tinyecs engine itself.
type Engine struct {
	nextComponentID uint64

	components   map[uint64]any
	componentMtx sync.RWMutex

	entities []ecsEntity

	links map[uint64]entityComponentLink
}

// AddComponents adds one or more component to the entity.
// This also updates the Engine's global component list.
func (e *Engine) AddComponents(entity ecsEntity, components ...any) {
	for _, component := range components {
		e.addComponent(entity, component)
	}
}

// DeleteComponent deletes a component.
func (e *Engine) DeleteComponent(component any) {
	for id, comp := range e.components {
		if comp == component {
			// Reset the link
			delete(e.links, id)

			//ent.components = append(ent.components[:i], ent.components[i+1:]...)
			e.deleteComponent(id)
		}
	}
}

// addComponent takes a slice of components and adds it to the engine and increments the nextComponentID variable.
func (e *Engine) addComponent(entity any, component any) uint64 {
	e.componentMtx.Lock()
	defer e.componentMtx.Unlock()

	id := e.nextComponentID
	e.components[id] = component

	// Set the link relationship.
	e.links[id] = entityComponentLink{
		entity:    entity,
		component: &component,
	}

	e.nextComponentID++
	return id
}

// deleteComponent is an internal function used to delete a component by id.
func (e *Engine) deleteComponent(id uint64) {
	e.componentMtx.Lock()
	defer e.componentMtx.Unlock()

	delete(e.components, id)
}

// GetComponents returns the components map held by the engine.
func (e *Engine) GetComponents() map[uint64]any {
	return e.components
}

// GetEntities returns a slice of entities held by the engine.
func (e *Engine) GetEntities() []ecsEntity {
	return e.entities
}

// AddEntity adds an entity to the engine.
func (e *Engine) AddEntity(entity ecsEntity) {
	e.entities = append(e.entities, entity)
}

// RemoveEntity takes in an entity instance and removes it from the engine.
// Note: This is pretty slow due to the use of reflect.DeepEqual.
func (e *Engine) RemoveEntity(entity ecsEntity) {
	for i, ent := range e.entities {
		// TODO: Replace DeepEqual since it is pretty slow.
		if reflect.DeepEqual(ent, entity) {
			e.entities = append(e.entities[:i], e.entities[i+1:]...)
			return
		}
	}
}

// NewEngine returns a prepared Engine instance ready for use.
// This should be the entry point for the tinyecs library.
func NewEngine() Engine {
	return Engine{
		links:      make(map[uint64]entityComponentLink),
		components: make(map[uint64]any),
	}
}

// Each is a generic function that iterates over the engine's components
// and runs the function argument with arguments being the id of the
// component along with the component instance of type T.
// Note that Each does not provide the actual entity used. Use EachEntity
// instead for this purpose.
//
// 		tinyecs.Each[Timer](&e, func(id uint64, obj Timer) {
//			obj.currentTime += 0.35
//			tinyecs.Set(&e, id, obj)
//		})
//
// The example above illustrates a basic use case where one updates a variable on a component, using the Set function.
func Each[T any](engine *Engine, f func(id uint64, component T)) uint64 {

	// Store a counter of objects touched which will be returned out of the function.
	var counter uint64

	// Iterate through all engine components.
	for idx, component := range engine.components {
		// Attempt to cast, and call the func on each of the components that can be successfully cast.
		if c, ok := component.(T); ok {
			counter++
			f(idx, c)
		}
	}
	return counter
}

// EachEntity is a generic function that iterates over the components belonging to entity E as component C.
//
// 		tinyecs.Each[MyEntity, Timer](&e, func(entity MyEntity, component Timer) {
//			log.Println("MyEntity: ", entity)
//		})
func EachEntity[E any, C any](engine *Engine, f func(entity E, component C)) uint64 {
	var counter uint64

	for idx, link := range engine.links {
		component := *link.component
		if _, ok := component.(C); ok {
			if e, entOk := link.entity.(E); entOk {
				counter++
				f(e, engine.components[idx].(C))
			}
		}
	}

	return counter
}

// Set takes in an engine instance and updates a component with the id specified.
func Set(engine *Engine, id uint64, component any) {
	engine.componentMtx.Lock()
	defer engine.componentMtx.Unlock()

	engine.components[id] = component
}

// ecsEntity is an internal type used to represent an entity.
type ecsEntity interface {
	GetComponents(engine *Engine) []uint64
}

// Entity is a collection of components.
// Consumers should extend this by embedding the struct.
type Entity struct{}

// GetComponents returns a list of component IDs associated with the entity.
func (ent Entity) GetComponents(engine *Engine) []uint64 {
	var components []uint64

	for i, link := range engine.links {
		if ent == link.entity {
			components = append(components, i)
		}
	}

	return components
}
