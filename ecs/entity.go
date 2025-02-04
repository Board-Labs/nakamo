package ecs

// Entity represents a unique identifier for an entity in the ECS system
type Entity uint64

// EntityManager handles entity creation and deletion
type EntityManager struct {
	currentID Entity
	recycled  []Entity
}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{
		currentID: 0,
		recycled:  make([]Entity, 0),
	}
}

// Create generates a new unique entity ID
func (em *EntityManager) Create() Entity {
	if len(em.recycled) > 0 {
		lastIdx := len(em.recycled) - 1
		entity := em.recycled[lastIdx]
		em.recycled = em.recycled[:lastIdx]
		return entity
	}

	em.currentID++
	return em.currentID
}

// Destroy marks an entity as destroyed and recycles its ID
func (em *EntityManager) Destroy(entity Entity) {
	em.recycled = append(em.recycled, entity)
}
