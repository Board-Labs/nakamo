package ecs

// System represents a system that processes entities with specific components
type System interface {
	Update(dt float64)
}

// World represents the game world containing all entities, components, and systems
type World struct {
	entities   *EntityManager
	components *ComponentRegistry
	systems    []System
}

// NewWorld creates a new game world
func NewWorld() *World {
	return &World{
		entities:   NewEntityManager(),
		components: NewComponentRegistry(),
		systems:    make([]System, 0),
	}
}

// CreateEntity creates a new entity in the world
func (w *World) CreateEntity() Entity {
	return w.entities.Create()
}

// DestroyEntity removes an entity and all its components from the world
func (w *World) DestroyEntity(entity Entity) {
	w.entities.Destroy(entity)
	for _, storage := range w.components.storages {
		if storage.Has(entity) {
			storage.Remove(entity)
		}
	}
}

// RegisterStorage registers a new component storage
func (w *World) RegisterStorage(storage ComponentStorage) ComponentID {
	return w.components.Register(storage)
}

// AddSystem adds a system to the world
func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

// Update updates all systems in the world
func (w *World) Update(dt float64) {
	for _, system := range w.systems {
		system.Update(dt)
	}
}

// Clear removes all entities and components from the world
func (w *World) Clear() {
	w.components.Clear()
	w.systems = make([]System, 0)
}

// GetStorage retrieves a component storage by its ID
func (w *World) GetStorage(id ComponentID) (ComponentStorage, bool) {
	return w.components.GetStorage(id)
}
