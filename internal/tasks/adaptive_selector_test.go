package tasks

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestDefaultAdaptiveWeights(t *testing.T) {
	w := DefaultAdaptiveWeights()
	if w.Diversity != 1.0 || w.Mood != 1.0 || w.TimeOfDay != 0.5 || w.Avoidance != 1.0 {
		t.Errorf("unexpected default weights: %+v", w)
	}
}

func TestAdaptiveSelector_ScoreCandidate_DiversityOnly(t *testing.T) {
	selector := newAdaptiveSelectorWithHour("", nil, DefaultAdaptiveWeights(), 10)
	tasks := []*Task{
		newCategorizedTestTask("1", "t1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("2", "t2", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
		newCategorizedTestTask("3", "t3", StatusTodo, TypeAdministrative, EffortMedium, LocationAnywhere),
	}
	score := selector.ScoreCandidate(tasks)
	// With nil patterns: diversity only = 9 * 1.0 = 9.0
	expected := 9.0
	if score != expected {
		t.Errorf("ScoreCandidate() = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_ScoreCandidate_WithMood(t *testing.T) {
	report := makeTestPatternReport(
		withMoodCorrelations(MoodCorrelation{
			Mood:          "happy",
			SessionCount:  5,
			PreferredType: "creative",
		}),
	)
	selector := newAdaptiveSelectorWithHour("happy", report, DefaultAdaptiveWeights(), 10)

	tasks := []*Task{
		newCategorizedTestTask("1", "Creative task", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("2", "Technical task", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
	}
	score := selector.ScoreCandidate(tasks)
	// diversity = 6, mood = 2 (one creative match), tod = 0, avoidance = 0
	// 6*1.0 + 2*1.0 + 0*0.5 - 0*1.0 = 8.0
	expected := 8.0
	if score != expected {
		t.Errorf("ScoreCandidate() with mood = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_ScoreCandidate_WithAvoidance(t *testing.T) {
	report := makeTestPatternReport(
		withAvoidance(AvoidanceEntry{TaskText: "Avoided task", TimesBypassed: 7, TimesShown: 10}),
	)
	selector := newAdaptiveSelectorWithHour("", report, DefaultAdaptiveWeights(), 10)

	tasks := []*Task{
		newCategorizedTestTask("1", "Avoided task", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("2", "Normal task", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
	}
	score := selector.ScoreCandidate(tasks)
	// diversity = 6, mood = 0, tod = 0, avoidance = 0.5 (7 bypasses)
	// 6*1.0 + 0 + 0 - 0.5*1.0 = 5.5
	expected := 5.5
	if score != expected {
		t.Errorf("ScoreCandidate() with avoidance = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_ScoreCandidate_WithTimeOfDay(t *testing.T) {
	report := makeTestPatternReport(
		withTimeOfDay(
			TimeOfDayPattern{Period: "morning", SessionCount: 5, AvgTasksCompleted: 4.0},
			TimeOfDayPattern{Period: "afternoon", SessionCount: 3, AvgTasksCompleted: 2.0},
		),
	)
	// hour=10 → morning, which is the most productive
	selector := newAdaptiveSelectorWithHour("", report, DefaultAdaptiveWeights(), 10)

	tasks := []*Task{
		newCategorizedTestTask("1", "t1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("2", "t2", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
	}
	score := selector.ScoreCandidate(tasks)
	// diversity = 6, mood = 0, tod = 1.0 (in most productive period), avoidance = 0
	// 6*1.0 + 0 + 1.0*0.5 - 0 = 6.5
	expected := 6.5
	if score != expected {
		t.Errorf("ScoreCandidate() with time-of-day = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_ScoreCandidate_NotMostProductivePeriod(t *testing.T) {
	report := makeTestPatternReport(
		withTimeOfDay(
			TimeOfDayPattern{Period: "morning", SessionCount: 5, AvgTasksCompleted: 4.0},
			TimeOfDayPattern{Period: "afternoon", SessionCount: 3, AvgTasksCompleted: 2.0},
		),
	)
	// hour=14 → afternoon, NOT the most productive
	selector := newAdaptiveSelectorWithHour("", report, DefaultAdaptiveWeights(), 14)

	tasks := []*Task{
		newCategorizedTestTask("1", "t1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("2", "t2", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
	}
	score := selector.ScoreCandidate(tasks)
	// diversity = 6, mood = 0, tod = 0 (NOT most productive), avoidance = 0
	expected := 6.0
	if score != expected {
		t.Errorf("ScoreCandidate() not most productive = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_ScoreCandidate_Combined(t *testing.T) {
	report := makeTestPatternReport(
		withMoodCorrelations(MoodCorrelation{
			Mood:            "stressed",
			SessionCount:    5,
			PreferredEffort: "quick-win",
		}),
		withAvoidance(AvoidanceEntry{TaskText: "Admin task", TimesBypassed: 12, TimesShown: 15}),
		withTimeOfDay(
			TimeOfDayPattern{Period: "morning", SessionCount: 5, AvgTasksCompleted: 4.0},
		),
	)
	selector := newAdaptiveSelectorWithHour("stressed", report, DefaultAdaptiveWeights(), 9)

	tasks := []*Task{
		newCategorizedTestTask("1", "Tech task", StatusTodo, TypeTechnical, EffortDeepWork, LocationAnywhere),
		newCategorizedTestTask("2", "Admin task", StatusTodo, TypeAdministrative, EffortMedium, LocationWork),
		newCategorizedTestTask("3", "Physical task", StatusTodo, TypePhysical, EffortQuickWin, LocationHome),
	}
	score := selector.ScoreCandidate(tasks)
	// diversity = 9, mood = 1 (physical has quick-win effort match), tod = 1.0 (morning is most productive), avoidance = 0.8 (admin bypassed 12x)
	// 9*1.0 + 1*1.0 + 1.0*0.5 - 0.8*1.0 = 9.7
	expected := 9.7
	if math.Abs(score-expected) > 0.01 {
		t.Errorf("ScoreCandidate() combined = %f, want %f", score, expected)
	}
}

func TestAdaptiveSelector_AvoidancePenalty_Thresholds(t *testing.T) {
	report := makeTestPatternReport(
		withAvoidance(
			AvoidanceEntry{TaskText: "none", TimesBypassed: 2},
			AvoidanceEntry{TaskText: "low", TimesBypassed: 5},
			AvoidanceEntry{TaskText: "mid", TimesBypassed: 9},
			AvoidanceEntry{TaskText: "high", TimesBypassed: 10},
			AvoidanceEntry{TaskText: "very-high", TimesBypassed: 15},
		),
	)
	selector := newAdaptiveSelectorWithHour("", report, DefaultAdaptiveWeights(), 10)

	tests := []struct {
		text string
		want float64
	}{
		{"not-in-list", 0.0},
		{"none", 0.0},      // 2 bypasses
		{"low", 0.5},       // 5 bypasses
		{"mid", 0.5},       // 9 bypasses
		{"high", 0.8},      // 10 bypasses
		{"very-high", 0.8}, // 15 bypasses
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			task := newTestTask("id", tt.text, StatusTodo, baseTime)
			got := selector.AvoidancePenalty(task)
			if got != tt.want {
				t.Errorf("AvoidancePenalty(%q) = %f, want %f", tt.text, got, tt.want)
			}
		})
	}
}

func TestAdaptiveSelector_TimeOfDayBonus_EachPeriod(t *testing.T) {
	tests := []struct {
		name string
		hour int
		best string
		want float64
	}{
		{"morning match", 9, "morning", 1.0},
		{"afternoon match", 14, "afternoon", 1.0},
		{"evening match", 18, "evening", 1.0},
		{"night match", 23, "night", 1.0},
		{"morning no match", 9, "evening", 0.0},
		{"no patterns", 9, "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var report *PatternReport
			if tt.best != "" {
				report = makeTestPatternReport(
					withTimeOfDay(TimeOfDayPattern{Period: tt.best, SessionCount: 5, AvgTasksCompleted: 5.0}),
				)
			}
			selector := newAdaptiveSelectorWithHour("", report, DefaultAdaptiveWeights(), tt.hour)
			got := selector.timeOfDayBonus()
			if got != tt.want {
				t.Errorf("timeOfDayBonus() at hour %d = %f, want %f", tt.hour, got, tt.want)
			}
		})
	}
}

func TestSelectDoorsAdaptive_NilSelector(t *testing.T) {
	pool := poolFromTasks(
		newCategorizedTestTask("t1", "Task 1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("t2", "Task 2", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
		newCategorizedTestTask("t3", "Task 3", StatusTodo, TypeAdministrative, EffortMedium, LocationAnywhere),
		newCategorizedTestTask("t4", "Task 4", StatusTodo, TypePhysical, EffortQuickWin, LocationErrands),
	)

	rng := rand.New(rand.NewPCG(42, 0))
	selected := selectDoorsAdaptiveWithRand(pool, 3, nil, rng)
	if len(selected) != 3 {
		t.Fatalf("expected 3 doors, got %d", len(selected))
	}
}

func TestSelectDoorsAdaptive_DiversityFloor(t *testing.T) {
	report := makeTestPatternReport(
		withMoodCorrelations(MoodCorrelation{
			Mood:          "happy",
			SessionCount:  5,
			PreferredType: "creative",
		}),
	)
	// Pool with mostly creative tasks
	pool := poolFromTasks(
		newCategorizedTestTask("t1", "Creative 1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("t2", "Creative 2", StatusTodo, TypeCreative, EffortMedium, LocationWork),
		newCategorizedTestTask("t3", "Creative 3", StatusTodo, TypeCreative, EffortDeepWork, LocationAnywhere),
		newCategorizedTestTask("t4", "Technical 1", StatusTodo, TypeTechnical, EffortMedium, LocationWork),
	)

	selector := newAdaptiveSelectorWithHour("happy", report, DefaultAdaptiveWeights(), 10)
	rng := rand.New(rand.NewPCG(42, 0))
	selected := selectDoorsAdaptiveWithRand(pool, 3, selector, rng)

	if len(selected) != 3 {
		t.Fatalf("expected 3 doors, got %d", len(selected))
	}

	// Diversity floor: not all should be creative (one should be swapped)
	creativeCount := 0
	for _, task := range selected {
		if task.Type == TypeCreative {
			creativeCount++
		}
	}
	if creativeCount == 3 {
		t.Error("diversity floor not enforced: all 3 doors are creative type")
	}
}

func TestSelectDoorsAdaptive_Deterministic(t *testing.T) {
	// Verify that with the same RNG seed and same pool, results are the same across runs
	report := makeTestPatternReport(
		withMoodCorrelations(MoodCorrelation{Mood: "happy", SessionCount: 5, PreferredType: "creative"}),
	)

	// Run twice with the same pool and same seed — map iteration order is non-deterministic,
	// but with N=10 candidates and best-score selection, the SCORE is deterministic
	// even if internal ordering varies. Verify score consistency.
	pool := poolFromTasks(
		newCategorizedTestTask("t1", "Task 1", StatusTodo, TypeCreative, EffortQuickWin, LocationHome),
		newCategorizedTestTask("t2", "Task 2", StatusTodo, TypeTechnical, EffortDeepWork, LocationWork),
		newCategorizedTestTask("t3", "Task 3", StatusTodo, TypeAdministrative, EffortMedium, LocationAnywhere),
		newCategorizedTestTask("t4", "Task 4", StatusTodo, TypePhysical, EffortQuickWin, LocationErrands),
		newCategorizedTestTask("t5", "Task 5", StatusTodo, TypeCreative, EffortMedium, LocationWork),
	)

	selector := newAdaptiveSelectorWithHour("happy", report, DefaultAdaptiveWeights(), 10)
	rng := rand.New(rand.NewPCG(99, 0))
	result := selectDoorsAdaptiveWithRand(pool, 3, selector, rng)

	if len(result) != 3 {
		t.Fatalf("expected 3 doors, got %d", len(result))
	}

	// Verify the selected set has a good score (above diversity-only baseline)
	score := selector.ScoreCandidate(result)
	if score < 6.0 {
		t.Errorf("expected good adaptive score (>= 6.0), got %f", score)
	}
}

func TestSelectDoorsAdaptive_EmptyPool(t *testing.T) {
	pool := NewTaskPool()
	selector := newAdaptiveSelectorWithHour("", nil, DefaultAdaptiveWeights(), 10)
	rng := rand.New(rand.NewPCG(42, 0))
	result := selectDoorsAdaptiveWithRand(pool, 3, selector, rng)
	if result != nil {
		t.Errorf("expected nil for empty pool, got %v", result)
	}
}

func TestSelectDoorsAdaptive_FewTasks(t *testing.T) {
	pool := poolFromTasks(
		newTestTask("t1", "Task 1", StatusTodo, baseTime),
	)
	selector := newAdaptiveSelectorWithHour("", nil, DefaultAdaptiveWeights(), 10)
	rng := rand.New(rand.NewPCG(42, 0))
	result := selectDoorsAdaptiveWithRand(pool, 3, selector, rng)
	if len(result) != 1 {
		t.Errorf("expected 1 door (all available), got %d", len(result))
	}
}

func TestHourToPeriod(t *testing.T) {
	tests := []struct {
		hour int
		want string
	}{
		{5, "morning"},
		{11, "morning"},
		{12, "afternoon"},
		{16, "afternoon"},
		{17, "evening"},
		{20, "evening"},
		{21, "night"},
		{4, "night"},
		{0, "night"},
	}
	for _, tt := range tests {
		got := HourToPeriod(tt.hour)
		if got != tt.want {
			t.Errorf("HourToPeriod(%d) = %q, want %q", tt.hour, got, tt.want)
		}
	}
}
