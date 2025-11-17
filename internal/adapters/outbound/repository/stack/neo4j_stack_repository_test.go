package stack

import (
	"context"
	domain "devconnectrelations/internal/domain/stack"
	"devconnectrelations/internal/tests"
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver neo4j.DriverWithContext
)

func TestMain(m *testing.M) {
	d, cleanup := tests.SetupNeo4j(&testing.T{})
	driver = d
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestNeo4jStackRepository_CreateStack(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testStack := &domain.Stack{Name: "TestStack"}
	repo := NewNeo4jStackRepository(driver)
	stack, err := repo.CreateStack(ctx, testStack)
	if err != nil {
		t.Fatalf("Failed to save stack: %v", err)
	}

	if stack.Name != testStack.Name {
		t.Errorf("Expected stack name %s, got %s", testStack.Name, stack.Name)
	}
}

func TestNeo4jStackRepository_GetStackByName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testStack := &domain.Stack{Name: "TestStackToGet"}
	repo := NewNeo4jStackRepository(driver)
	_, err := repo.CreateStack(ctx, testStack)
	if err != nil {
		t.Fatalf("Failed to save stack: %v", err)
	}
	fetchedStack, err := repo.GetStackByName(ctx, testStack.Name)
	if err != nil {
		t.Fatalf("Failed to get stack: %v", err)
	}
	if fetchedStack.Name != testStack.Name {
		t.Errorf("Expected stack name %s, got %s", testStack.Name, fetchedStack.Name)
	}
}

func TestNeo4jStackRepository_DeleteStack(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testStack := &domain.Stack{Name: "TestStackToDelete"}
	repo := NewNeo4jStackRepository(driver)
	_, err := repo.CreateStack(ctx, testStack)
	if err != nil {
		t.Fatalf("Failed to save stack: %v", err)
	}
	err = repo.DeleteStack(ctx, testStack.Name)
	if err != nil {
		t.Fatalf("Failed to delete stack: %v", err)
	}
	fetchedStack, err := repo.GetStackByName(ctx, testStack.Name)
	if err != nil {
		t.Fatalf("Failed to get stack: %v", err)
	}
	if fetchedStack.Name != "" {
		t.Errorf("Expected stack to be deleted, but found: %s", fetchedStack.Name)
	}
}
