package core

import (
	"fmt"
	"testing"
)

func BenchmarkTaskPoolAddTask(b *testing.B) {
	sizes := []int{10, 100, 500}
	for _, n := range sizes {
		b.Run(fmt.Sprintf("pool=%d", n), func(b *testing.B) {
			for range b.N {
				pool := NewTaskPool()
				for i := range n {
					pool.AddTask(&Task{
						ID:     fmt.Sprintf("t-%d", i),
						Text:   fmt.Sprintf("Task %d", i),
						Status: StatusTodo,
					})
				}
			}
		})
	}
}

func BenchmarkTaskPoolGetAllTasks(b *testing.B) {
	sizes := []int{10, 100, 500}
	for _, n := range sizes {
		pool := benchPool(n)
		b.Run(fmt.Sprintf("pool=%d", n), func(b *testing.B) {
			for range b.N {
				pool.GetAllTasks()
			}
		})
	}
}

func BenchmarkTaskPoolGetTasksByStatus(b *testing.B) {
	sizes := []int{10, 100, 500}
	for _, n := range sizes {
		pool := benchPool(n)
		b.Run(fmt.Sprintf("pool=%d", n), func(b *testing.B) {
			for range b.N {
				pool.GetTasksByStatus(StatusTodo)
			}
		})
	}
}

func BenchmarkTaskPoolGetAvailableForDoors(b *testing.B) {
	sizes := []int{10, 100, 500}
	for _, n := range sizes {
		pool := benchPool(n)
		b.Run(fmt.Sprintf("pool=%d", n), func(b *testing.B) {
			for range b.N {
				pool.GetAvailableForDoors()
			}
		})
	}
}

func BenchmarkTaskPoolGetTask(b *testing.B) {
	sizes := []int{10, 100, 500}
	for _, n := range sizes {
		pool := benchPool(n)
		targetID := fmt.Sprintf("bench-%d", n/2)
		b.Run(fmt.Sprintf("pool=%d", n), func(b *testing.B) {
			for range b.N {
				pool.GetTask(targetID)
			}
		})
	}
}
