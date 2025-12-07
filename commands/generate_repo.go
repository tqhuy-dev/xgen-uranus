package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var repoName string

var generateRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Generate a new repository",
	Long: `Generate a new repository with the specified name.

This command creates a new repository scaffold with all necessary files
and directory structure for data access layer.

Examples:
  uranus generate repo --name simple_repo
  uranus generate repo --name path/simple_repo
  uranus generate repo --name internal/repository/user_repo`,
	Run: func(cmd *cobra.Command, args []string) {
		if repoName == "" {
			fmt.Println("Error: --name flag is required")
			cmd.Help()
			os.Exit(1)
		}

		if err := generateRepo(repoName); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating repo: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	generateRepoCmd.Flags().StringVarP(&repoName, "name", "n", "", "Name of the repository to generate (required)")
	generateRepoCmd.MarkFlagRequired("name")
}

func generateRepo(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("repo name cannot be empty")
	}

	// Get the directory path and repo name
	dir := filepath.Dir(name)
	baseName := filepath.Base(name)

	// Determine the full path for the repo
	repoDir := name
	if dir == "." {
		repoDir = baseName
	}

	fmt.Printf("🗄️  Generating repository: %s\n", baseName)

	// Create the repo directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", repoDir, err)
	}
	fmt.Printf("  📁 Created: %s\n", repoDir)

	// Create repository interface file
	interfaceContent := fmt.Sprintf(`package %s

import (
	"context"
)

// Entity represents the data model for %s
type Entity struct {
	ID        string
	CreatedAt int64
	UpdatedAt int64
}

// Repository defines the interface for %s data access
type Repository interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *Entity) error
	
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (*Entity, error)
	
	// Update updates an existing entity
	Update(ctx context.Context, entity *Entity) error
	
	// Delete removes an entity by its ID
	Delete(ctx context.Context, id string) error
	
	// List retrieves all entities with pagination
	List(ctx context.Context, offset, limit int) ([]*Entity, error)
}
`, baseName, baseName, baseName)

	interfacePath := filepath.Join(repoDir, "repository.go")
	if err := os.WriteFile(interfacePath, []byte(interfaceContent), 0644); err != nil {
		return fmt.Errorf("failed to create repository.go: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", interfacePath)

	// Create repository implementation file
	implContent := fmt.Sprintf(`package %s

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryRepository is an in-memory implementation of Repository
type InMemoryRepository struct {
	mu       sync.RWMutex
	entities map[string]*Entity
}

// NewInMemoryRepository creates a new InMemoryRepository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		entities: make(map[string]*Entity),
	}
}

// Create creates a new entity
func (r *InMemoryRepository) Create(ctx context.Context, entity *Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entities[entity.ID]; exists {
		return fmt.Errorf("entity with ID %%s already exists", entity.ID)
	}

	r.entities[entity.ID] = entity
	return nil
}

// GetByID retrieves an entity by its ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, exists := r.entities[id]
	if !exists {
		return nil, fmt.Errorf("entity with ID %%s not found", id)
	}

	return entity, nil
}

// Update updates an existing entity
func (r *InMemoryRepository) Update(ctx context.Context, entity *Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entities[entity.ID]; !exists {
		return fmt.Errorf("entity with ID %%s not found", entity.ID)
	}

	r.entities[entity.ID] = entity
	return nil
}

// Delete removes an entity by its ID
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entities[id]; !exists {
		return fmt.Errorf("entity with ID %%s not found", id)
	}

	delete(r.entities, id)
	return nil
}

// List retrieves all entities with pagination
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]*Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entities := make([]*Entity, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity)
	}

	// Apply pagination
	if offset >= len(entities) {
		return []*Entity{}, nil
	}

	end := offset + limit
	if end > len(entities) {
		end = len(entities)
	}

	return entities[offset:end], nil
}

// Ensure InMemoryRepository implements Repository interface
var _ Repository = (*InMemoryRepository)(nil)
`, baseName)

	implPath := filepath.Join(repoDir, "repository_impl.go")
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		return fmt.Errorf("failed to create repository_impl.go: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", implPath)

	// Create test file
	testContent := fmt.Sprintf(`package %s

import (
	"context"
	"testing"
)

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	entity := &Entity{
		ID:        "test-1",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	err := repo.Create(ctx, entity)
	if err != nil {
		t.Errorf("Create() error = %%v, want nil", err)
	}

	// Test duplicate creation
	err = repo.Create(ctx, entity)
	if err == nil {
		t.Error("Create() expected error for duplicate, got nil")
	}
}

func TestInMemoryRepository_GetByID(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	entity := &Entity{
		ID:        "test-1",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	_ = repo.Create(ctx, entity)

	got, err := repo.GetByID(ctx, "test-1")
	if err != nil {
		t.Errorf("GetByID() error = %%v, want nil", err)
	}
	if got.ID != entity.ID {
		t.Errorf("GetByID() got ID = %%v, want %%v", got.ID, entity.ID)
	}

	// Test not found
	_, err = repo.GetByID(ctx, "not-found")
	if err == nil {
		t.Error("GetByID() expected error for not found, got nil")
	}
}

func TestInMemoryRepository_Update(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	entity := &Entity{
		ID:        "test-1",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	_ = repo.Create(ctx, entity)

	entity.UpdatedAt = 9999999999
	err := repo.Update(ctx, entity)
	if err != nil {
		t.Errorf("Update() error = %%v, want nil", err)
	}

	got, _ := repo.GetByID(ctx, "test-1")
	if got.UpdatedAt != 9999999999 {
		t.Errorf("Update() UpdatedAt = %%v, want 9999999999", got.UpdatedAt)
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	entity := &Entity{
		ID:        "test-1",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	_ = repo.Create(ctx, entity)

	err := repo.Delete(ctx, "test-1")
	if err != nil {
		t.Errorf("Delete() error = %%v, want nil", err)
	}

	_, err = repo.GetByID(ctx, "test-1")
	if err == nil {
		t.Error("GetByID() expected error after delete, got nil")
	}
}

func TestInMemoryRepository_List(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		entity := &Entity{
			ID:        "test-" + string(rune('0'+i)),
			CreatedAt: 1234567890,
			UpdatedAt: 1234567890,
		}
		_ = repo.Create(ctx, entity)
	}

	list, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Errorf("List() error = %%v, want nil", err)
	}
	if len(list) != 5 {
		t.Errorf("List() len = %%v, want 5", len(list))
	}

	// Test pagination
	list, err = repo.List(ctx, 2, 2)
	if err != nil {
		t.Errorf("List() error = %%v, want nil", err)
	}
	if len(list) != 2 {
		t.Errorf("List() with pagination len = %%v, want 2", len(list))
	}
}
`, baseName)

	testPath := filepath.Join(repoDir, "repository_test.go")
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to create repository_test.go: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", testPath)

	fmt.Printf("\n✅ Repository '%s' generated successfully!\n", baseName)
	fmt.Printf("\nFiles created:\n")
	fmt.Printf("  - %s (Repository interface)\n", interfacePath)
	fmt.Printf("  - %s (In-memory implementation)\n", implPath)
	fmt.Printf("  - %s (Unit tests)\n", testPath)

	return nil
}


