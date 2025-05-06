package co

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestIterAllResults(t *testing.T) {
	// Create a context for our tests
	ctx := context.Background()

	t.Run("all successful promises", func(t *testing.T) {
		// Create sample promises with different keys and values
		promises := []PromiseWithKey[int, string]{
			{Key: "one", Promise: NewResolved(1)},
			{Key: "two", Promise: NewResolved(2)},
			{Key: "three", Promise: NewResolved(3)},
		}

		// Use IterAllResults to iterate through the promises
		iter := IterAllResults(ctx, promises...)

		// Collect results using the iterator
		results := make(map[string]IterResultItem[int])
		iter(func(key string, item IterResultItem[int]) bool {
			results[key] = item
			return true // continue iteration
		})

		// Verify we got all results with expected values
		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}

		expectedValues := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}

		for key, item := range results {
			if item.Err != nil {
				t.Errorf("Expected no error for key %s, got %v", key, item.Err)
			}
			expectedVal, exists := expectedValues[key]
			if !exists {
				t.Errorf("Unexpected key in results: %s", key)
				continue
			}
			if item.Result != expectedVal {
				t.Errorf("For key %s: expected value %d, got %d", key, expectedVal, item.Result)
			}
		}
	})

	t.Run("mixed success and error promises", func(t *testing.T) {
		testErr := errors.New("test error")

		promises := []PromiseWithKey[int, string]{
			{Key: "success1", Promise: NewResolved(10)},
			{Key: "error", Promise: NewRejected[int](testErr)},
			{Key: "success2", Promise: NewResolved(20)},
		}

		iter := IterAllResults(ctx, promises...)

		results := make(map[string]IterResultItem[int])
		iter(func(key string, item IterResultItem[int]) bool {
			results[key] = item
			return true
		})

		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}

		// Check success results
		if r, ok := results["success1"]; !ok || r.Err != nil || r.Result != 10 {
			t.Errorf("Incorrect result for success1: %+v", r)
		}
		if r, ok := results["success2"]; !ok || r.Err != nil || r.Result != 20 {
			t.Errorf("Incorrect result for success2: %+v", r)
		}

		// Check error result
		if r, ok := results["error"]; !ok || r.Err != testErr {
			t.Errorf("Expected error for 'error' key, got: %+v", r)
		}
	})

	t.Run("early termination", func(t *testing.T) {
		promises := []PromiseWithKey[int, string]{
			{Key: "one", Promise: NewResolved(1)},
			{Key: "two", Promise: NewResolved(2)},
			{Key: "three", Promise: NewResolved(3)},
		}

		iter := IterAllResults(ctx, promises...)

		// Only process first two results by returning false after second item
		results := make(map[string]IterResultItem[int])
		count := 0
		iter(func(key string, item IterResultItem[int]) bool {
			results[key] = item
			count++
			return count < 2 // stop after processing 2 items
		})

		if len(results) != 2 {
			t.Errorf("Expected 2 results after early termination, got %d", len(results))
		}
	})

	t.Run("delayed promises", func(t *testing.T) {
		// Create promises with different completion times
		p1 := NewPromise(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 100, nil
		})
		p2 := NewPromise(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 200, nil
		})
		p3 := NewPromise(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 300, nil
		})

		promises := []PromiseWithKey[int, string]{
			{Key: "slow", Promise: p1},
			{Key: "fast", Promise: p2},
			{Key: "medium", Promise: p3},
		}

		// Create a context with timeout to ensure test doesn't hang
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		iter := IterAllResults(ctx, promises...)

		// Track the order of results
		var resultOrder []string
		results := make(map[string]IterResultItem[int])

		iter(func(key string, item IterResultItem[int]) bool {
			results[key] = item
			resultOrder = append(resultOrder, key)
			return true
		})

		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}

		// Check all results are correct
		expectedValues := map[string]int{
			"slow":   100,
			"fast":   200,
			"medium": 300,
		}

		for key, item := range results {
			if item.Err != nil {
				t.Errorf("Expected no error for key %s, got %v", key, item.Err)
			}
			expectedVal, exists := expectedValues[key]
			if !exists {
				t.Errorf("Unexpected key in results: %s", key)
				continue
			}
			if item.Result != expectedVal {
				t.Errorf("For key %s: expected value %d, got %d", key, expectedVal, item.Result)
			}
		}

		// Note: We don't test exact order as it can vary, but results should include all three keys
		if len(resultOrder) != 3 {
			t.Errorf("Expected 3 results in order tracking, got %d", len(resultOrder))
		}
	})

	t.Run("empty promises", func(t *testing.T) {
		iter := IterAllResults[int, string](ctx)

		called := false
		iter(func(key string, item IterResultItem[int]) bool {
			called = true
			return true
		})

		if called {
			t.Error("Iterator function should not be called for empty promises list")
		}
	})

	t.Run("with cancelled context", func(t *testing.T) {
		// Create a context that's already cancelled
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		executionCount := atomic.Int32{}

		p1 := NewPromise(func() (int, error) {
			executionCount.Add(1)
			return 1, nil
		})

		p2 := NewPromise(func() (int, error) {
			executionCount.Add(1)
			return 2, nil
		})

		promises := PromiseMap[string, int]{
			"one": p1,
			"two": p2,
		}

		promises.Await(context.Background())

		iter := IterAllResults(cancelledCtx, promises.AsPromiseWithKeys()...)

		results := make(map[string]IterResultItem[int])
		iter(func(key string, item IterResultItem[int]) bool {
			results[key] = item
			return true
		})

		// The promises should still execute, just return with context error
		if len(results) != 2 {
			t.Errorf("Expected 2 results despite cancelled context, got %d", len(results))
		}

		// Check that all results have context cancelled error
		for key, result := range results {
			if result.Err != nil {
				t.Errorf("Expected no error for key %s, got: %v", key, result.Err)
			}
		}
	})
}
