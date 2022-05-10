# tinyecs
Entity Component System written in Golang using generics.

## Basic usage
```go
// Initialize a new instance of the ECS engine.
e := tinyecs.NewEngine()

// Define and set up a new testEntity.
entity := testEntity{
	tinyecs.Entity
}

// Add some components along with passing in the engine instance.
entity.AddComponents(
    &e,
	
    velocity{},
    playerData{name: "test", health: 100.0},
)

// Add the entity to the ECS engine.
e.AddEntity(entity)

// This use of the Each function iterates over all playerdata components
// and prints the name of each of them.
tinyecs.Each[playerData](&e, func(id uint64, obj playerData) {
	log.Println(obj.name)
})
```
