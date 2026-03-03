package tasks

import (
	"testing"
	"time"
)

func TestTaskPool_AddAndGet(t *testing.T) {
	pool := NewTaskPool()
	task := NewTask("Test task")
	pool.AddTask(task)

	got := pool.GetTask(task.ID)
	if got == nil {
		t.Fatal("Expected to get task back")
	}
	if got.Text != task.Text {
		t.Errorf("Expected %q, got %q", task.Text, got.Text)
	}
}

func TestTaskPool_RemoveTask(t *testing.T) {
	pool := NewTaskPool()
	task := NewTask("Test")
	pool.AddTask(task)
	pool.RemoveTask(task.ID)

	if pool.GetTask(task.ID) != nil {
		t.Error("Expected task to be removed")
	}
	if pool.Count() != 0 {
		t.Errorf("Expected 0 tasks, got %d", pool.Count())
	}
}

func TestTaskPool_GetTasksByStatus(t *testing.T) {
	pool := NewTaskPool()
	t1 := NewTask("Todo task")
	t2 := NewTask("Blocked task")
	_ = t2.UpdateStatus(StatusBlocked)
	pool.AddTask(t1)
	pool.AddTask(t2)

	todos := pool.GetTasksByStatus(StatusTodo)
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo task, got %d", len(todos))
	}

	blocked := pool.GetTasksByStatus(StatusBlocked)
	if len(blocked) != 1 {
		t.Errorf("Expected 1 blocked task, got %d", len(blocked))
	}
}

func TestTaskPool_GetAvailableForDoors(t *testing.T) {
	pool := NewTaskPool()
	for i := 0; i < 5; i++ {
		pool.AddTask(NewTask("Task"))
	}

	available := pool.GetAvailableForDoors()
	if len(available) != 5 {
		t.Errorf("Expected 5 available tasks, got %d", len(available))
	}

	// Complete one task
	allTasks := pool.GetAllTasks()
	_ = allTasks[0].UpdateStatus(StatusComplete)
	pool.UpdateTask(allTasks[0])

	available = pool.GetAvailableForDoors()
	if len(available) != 4 {
		t.Errorf("Expected 4 available tasks after completing one, got %d", len(available))
	}
}

func TestTaskPool_RecentlyShown(t *testing.T) {
	pool := NewTaskPool()
	task := NewTask("Test")
	pool.AddTask(task)

	if pool.IsRecentlyShown(task.ID) {
		t.Error("Task should not be recently shown initially")
	}

	pool.MarkRecentlyShown(task.ID)
	if !pool.IsRecentlyShown(task.ID) {
		t.Error("Task should be recently shown after marking")
	}
}

func TestTaskPool_GetAvailableForDoors_FewTasks(t *testing.T) {
	pool := NewTaskPool()
	t1 := NewTask("Only task")
	pool.AddTask(t1)
	pool.MarkRecentlyShown(t1.ID)

	// With < 3 tasks, should include recently shown
	available := pool.GetAvailableForDoors()
	if len(available) != 1 {
		t.Errorf("Expected 1 available task (including recently shown), got %d", len(available))
	}
}

func TestGetAvailableForDoors_ExcludesDeferred(t *testing.T) {
	pool := NewTaskPool()
	t1 := NewTask("Normal task")
	t2 := NewTask("Deferred task")
	future := time.Now().Add(24 * time.Hour)
	t2.DeferredUntil = &future
	pool.AddTask(t1)
	pool.AddTask(t2)

	available := pool.GetAvailableForDoors()
	for _, task := range available {
		if task.ID == t2.ID {
			t.Error("deferred task should be excluded from available doors")
		}
	}
	if len(available) != 1 {
		t.Errorf("expected 1 available task, got %d", len(available))
	}
}

func TestGetAvailableForDoors_IncludesExpiredDeferred(t *testing.T) {
	pool := NewTaskPool()
	t1 := NewTask("Expired deferred task")
	past := time.Now().Add(-24 * time.Hour)
	t1.DeferredUntil = &past
	pool.AddTask(t1)

	available := pool.GetAvailableForDoors()
	if len(available) != 1 {
		t.Errorf("expected 1 available task (expired deferral), got %d", len(available))
	}
}

func TestGetAvailableForDoors_DeferredExactlyNow(t *testing.T) {
	pool := NewTaskPool()
	t1 := NewTask("Borderline task")
	now := time.Now()
	t1.DeferredUntil = &now
	pool.AddTask(t1)

	// DeferredUntil.After(time.Now()) should be false when equal → task included
	available := pool.GetAvailableForDoors()
	if len(available) != 1 {
		t.Errorf("expected 1 available task (deferral expired at now), got %d", len(available))
	}
}
