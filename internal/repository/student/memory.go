package repository

import (
	"context"
	"sync"

	"todo/internal/domain/student"
)

// MemoryRepository is an in-memory implementation of Repository for testing.
type MemoryRepository struct {
	mu       sync.RWMutex
	students map[string]*student.Student
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		students: make(map[string]*student.Student),
	}
}

// Save saves a new student record.
func (r *MemoryRepository) Save(ctx context.Context, s *student.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[s.StudentNumber]; exists {
		return student.NewStudentNumberAlreadyExistsError()
	}

	r.students[s.StudentNumber] = s
	return nil
}

// FindByStudentNumber retrieves a student by student number.
func (r *MemoryRepository) FindByStudentNumber(ctx context.Context, studentNumber string) (*student.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, exists := r.students[studentNumber]
	if !exists {
		return nil, student.NewStudentNotFoundError()
	}

	return s, nil
}

// FindAll retrieves all student records.
func (r *MemoryRepository) FindAll(ctx context.Context) ([]*student.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	students := make([]*student.Student, 0, len(r.students))
	for _, s := range r.students {
		students = append(students, s)
	}

	return students, nil
}

// Update updates an existing student record.
func (r *MemoryRepository) Update(ctx context.Context, s *student.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[s.StudentNumber]; !exists {
		return student.NewStudentNotFoundError()
	}

	r.students[s.StudentNumber] = s
	return nil
}

// Delete deletes a student record by student number.
func (r *MemoryRepository) Delete(ctx context.Context, studentNumber string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.students[studentNumber]; !exists {
		return student.NewStudentNotFoundError()
	}

	delete(r.students, studentNumber)
	return nil
}

// ExistsByStudentNumber checks if a student number exists.
func (r *MemoryRepository) ExistsByStudentNumber(ctx context.Context, studentNumber string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.students[studentNumber]
	return exists, nil
}
