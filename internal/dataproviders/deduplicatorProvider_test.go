package dataproviders

import (
	"reflect"
	"testing"
)

// Mock DataProvider for testing
type mockDataProvider struct {
	items []Item
	err   error
}

func (m *mockDataProvider) GetItems() ([]Item, error) {
	return m.items, m.err
}

func TestDeduplicatorProvider_GetItems_NoDuplicates(t *testing.T) {
	dirProvider := &mockDataProvider{
		items: []Item{
			{Id: "dir1", Display: "[ ] dir1"},
			{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
				{Id: "dir3", Display: "[ ] dir3"},
			}},
		},
	}
	muxProvider := &mockDataProvider{
		items: []Item{
			{Id: "mux1", Display: "mux1"},
		},
	}

	deduplicator := NewDeduplicatorProvider(dirProvider, muxProvider)
	result, err := deduplicator.GetItems()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Item{
		{Id: "dir1", Display: "[ ] dir1"},
		{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
			{Id: "dir3", Display: "[ ] dir3"},
		}},
		{Id: "mux1", Display: "mux1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDeduplicatorProvider_GetItems_WithDuplicates(t *testing.T) {
	dirProvider := &mockDataProvider{
		items: []Item{
			{Id: "dir1", Display: "[ ] dir1"},
			{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
				{Id: "dir3", Display: "[ ] dir3"},
			}},
		},
	}
	muxProvider := &mockDataProvider{
		items: []Item{
			{Id: "mux1", Display: "mux1"},
			{Id: "dir1", Display: "dir1"}, // Duplicate
		},
	}

	deduplicator := NewDeduplicatorProvider(dirProvider, muxProvider)
	result, err := deduplicator.GetItems()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Item{
		{Id: "dir1", Display: "[ ] dir1"},
		{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
			{Id: "dir3", Display: "[ ] dir3"},
		}},
		{Id: "mux1", Display: "mux1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDeduplicatorProvider_GetItems_MarkDuplicates(t *testing.T) {
	dirProvider := &mockDataProvider{
		items: []Item{
			{Id: "dir1", Display: "[ ] dir1"}, // Should be marked
			{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
				{Id: "dir3", Display: "[ ] dir3"}, // Should be marked
			}},
			{Id: "dir4", Display: "[ ] dir4"}, // Not in mux, should not be marked
		},
	}
	muxProvider := &mockDataProvider{
		items: []Item{
			{Id: "mux1", Display: "mux1"},
			{Id: "dir1", Display: "dir1"}, // Duplicate
			{Id: "dir3", Display: "dir3"}, // Duplicate
		},
	}

	deduplicator := NewDeduplicatorProvider(dirProvider, muxProvider).WithMarkDuplicates(true)
	result, err := deduplicator.GetItems()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Item{
		{Id: "dir1", Display: "[x] dir1"},
		{Id: "dir2", Display: "[ ] dir2", SubItems: []Item{
			{Id: "dir3", Display: "[x] dir3"},
		}},
		{Id: "dir4", Display: "[ ] dir4"},
		{Id: "mux1", Display: "mux1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDeduplicatorProvider_GetItems_MarkDuplicates_NoDuplicates(t *testing.T) {
	dirProvider := &mockDataProvider{
		items: []Item{
			{Id: "dir1", Display: "[ ] dir1"},
		},
	}
	muxProvider := &mockDataProvider{
		items: []Item{
			{Id: "mux1", Display: "mux1"},
		},
	}

	deduplicator := NewDeduplicatorProvider(dirProvider, muxProvider).WithMarkDuplicates(true)
	result, err := deduplicator.GetItems()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Item{
		{Id: "dir1", Display: "[ ] dir1"},
		{Id: "mux1", Display: "mux1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
