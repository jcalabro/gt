package gt_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jcalabro/gt"
)

func TestConcurrent(t *testing.T) {
	// Basic: double each item, verify order preserved.
	{
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		results, err := gt.Concurrent(context.Background(), items, func(n int) (int, error) {
			return n * 2, nil
		})
		require.NoError(t, err)
		require.Equal(t, []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, results)
	}

	// Empty input.
	{
		results, err := gt.Concurrent(context.Background(), []int{}, func(n int) (string, error) {
			return "", nil
		})
		require.NoError(t, err)
		require.Empty(t, results)
	}

	// Nil input.
	{
		results, err := gt.Concurrent(context.Background(), nil, func(n int) (string, error) {
			return "", nil
		})
		require.NoError(t, err)
		require.Empty(t, results)
	}

	// Single error returns error and nil results.
	{
		items := []int{1, 2, 3}
		results, err := gt.Concurrent(context.Background(), items, func(n int) (int, error) {
			if n == 2 {
				return 0, fmt.Errorf("bad: %d", n)
			}
			return n, nil
		})
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "bad: 2")
	}

	// Multiple errors are joined.
	{
		items := []int{1, 2, 3, 4}
		results, err := gt.Concurrent(context.Background(), items, func(n int) (int, error) {
			if n%2 == 0 {
				return 0, fmt.Errorf("even: %d", n)
			}
			return n, nil
		})
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "even: 2")
		require.Contains(t, err.Error(), "even: 4")
	}

	// All items error.
	{
		items := []int{1, 2, 3}
		results, err := gt.Concurrent(context.Background(), items, func(n int) (int, error) {
			return 0, fmt.Errorf("fail: %d", n)
		})
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "fail: 1")
		require.Contains(t, err.Error(), "fail: 2")
		require.Contains(t, err.Error(), "fail: 3")
	}

	// Panic in fn is recovered and surfaced as an error.
	{
		items := []int{1, 2, 3}
		results, err := gt.Concurrent(context.Background(), items, func(n int) (int, error) {
			if n == 2 {
				panic("boom")
			}
			return n, nil
		})
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "boom")
	}
}

func TestConcurrentN(t *testing.T) {
	// Custom worker count with errors.
	{
		items := make([]int, 100)
		for i := range items {
			items[i] = i
		}

		results, err := gt.ConcurrentN(context.Background(), items, 2, func(n int) (int, error) {
			if n == 50 {
				return 0, fmt.Errorf("mid-fail")
			}
			return n * 3, nil
		})
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "mid-fail")
	}

	// Custom worker count, all succeed.
	{
		items := make([]int, 100)
		for i := range items {
			items[i] = i
		}

		results, err := gt.ConcurrentN(context.Background(), items, 2, func(n int) (int, error) {
			return n * 3, nil
		})
		require.NoError(t, err)
		require.Len(t, results, 100)
		for i, r := range results {
			require.Equal(t, i*3, r)
		}
	}

	// Single worker (sequential execution).
	{
		items := []int{10, 20, 30}
		results, err := gt.ConcurrentN(context.Background(), items, 1, func(n int) (int, error) {
			return n + 1, nil
		})
		require.NoError(t, err)
		require.Equal(t, []int{11, 21, 31}, results)
	}

	// Context cancelled during dispatch phase with many items and 1 worker.
	{
		ctx, cancel := context.WithCancel(context.Background())

		items := make([]int, 100)
		for i := range items {
			items[i] = i
		}

		results, err := gt.ConcurrentN(ctx, items, 1, func(n int) (int, error) {
			if n == 2 {
				cancel()
			}
			return n, nil
		})
		require.ErrorIs(t, err, context.Canceled)
		require.Nil(t, results)
	}

	// Zero workers returns error.
	{
		results, err := gt.ConcurrentN(context.Background(), []int{1}, 0, func(n int) (int, error) {
			return n, nil
		})
		require.Error(t, err)
		require.Nil(t, results)
	}
}
