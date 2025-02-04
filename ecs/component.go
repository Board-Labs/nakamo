package ecs

import (
	"fmt"
	"reflect"
)

// ComponentID represents a unique identifier for a component type
type ComponentID uint32

// ComponentStorage interface defines the methods required for component storage
type ComponentStorage interface {
	Has(entity Entity) bool
	Remove(entity Entity)
	Clear()
}

// Storage is a generic component storage implementation
type Storage[T any] struct {
	components map[Entity]T
}

// NewStorage creates a new component storage
func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		components: make(map[Entity]T),
	}
}

// Add adds a component to an entity
func (s *Storage[T]) Add(entity Entity, component T) {
	s.components[entity] = component
}

// Get retrieves a component for an entity
func (s *Storage[T]) Get(entity Entity) (T, error) {
	if component, ok := s.components[entity]; ok {
		return component, nil
	}
	var zero T
	return zero, fmt.Errorf("entity %d does not have component %T", entity, zero)
}

// GetAll returns all entities that have this component
func (s *Storage[T]) GetAll() map[Entity]T {
	return s.components
}

// Has checks if an entity has a component
func (s *Storage[T]) Has(entity Entity) bool {
	_, ok := s.components[entity]
	return ok
}

// Remove removes a component from an entity
func (s *Storage[T]) Remove(entity Entity) {
	delete(s.components, entity)
}

// Clear removes all components
func (s *Storage[T]) Clear() {
	s.components = make(map[Entity]T)
}

// ComponentRegistry manages component type registration and storage creation
type ComponentRegistry struct {
	nextID     ComponentID
	components map[reflect.Type]ComponentID
	storages   map[ComponentID]ComponentStorage
}

// NewComponentRegistry creates a new component registry
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		nextID:     0,
		components: make(map[reflect.Type]ComponentID),
		storages:   make(map[ComponentID]ComponentStorage),
	}
}

// Register registers a new component type and returns its ID
func (r *ComponentRegistry) Register(storage ComponentStorage) ComponentID {
	storageType := reflect.TypeOf(storage).Elem()
	if id, exists := r.components[storageType]; exists {
		return id
	}

	id := r.nextID
	r.nextID++
	r.components[storageType] = id
	r.storages[id] = storage
	return id
}

// GetStorage retrieves the storage for a component type
func (r *ComponentRegistry) GetStorage(id ComponentID) (ComponentStorage, bool) {
	storage, ok := r.storages[id]
	return storage, ok
}

// Clear removes all components from all storages
func (r *ComponentRegistry) Clear() {
	for _, storage := range r.storages {
		storage.Clear()
	}
}
