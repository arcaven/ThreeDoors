package tasks

import (
	"math/rand/v2"
	"strings"
	"time"
)

// AdaptiveWeights controls the relative influence of each signal in door selection.
type AdaptiveWeights struct {
	Diversity float64 // default 1.0
	Mood      float64 // default 1.0
	TimeOfDay float64 // default 0.5
	Avoidance float64 // default 1.0
}

// DefaultAdaptiveWeights returns the default weight configuration.
func DefaultAdaptiveWeights() AdaptiveWeights {
	return AdaptiveWeights{
		Diversity: 1.0,
		Mood:      1.0,
		TimeOfDay: 0.5,
		Avoidance: 1.0,
	}
}

// AdaptiveSelector scores candidate door sets using multiple signals.
type AdaptiveSelector struct {
	weights              AdaptiveWeights
	currentMood          string
	patterns             *PatternReport
	currentHour          int
	avoidanceMap         map[string]int // task text → bypass count
	preferredType        TaskType
	preferredEffort      TaskEffort
	mostProductivePeriod string
}

// NewAdaptiveSelector creates a new AdaptiveSelector.
// currentHour is set from time.Now() internally.
func NewAdaptiveSelector(mood string, patterns *PatternReport, weights AdaptiveWeights) *AdaptiveSelector {
	return newAdaptiveSelectorWithHour(mood, patterns, weights, time.Now().Hour())
}

// newAdaptiveSelectorWithHour creates an AdaptiveSelector with an explicit hour for testing.
func newAdaptiveSelectorWithHour(mood string, patterns *PatternReport, weights AdaptiveWeights, hour int) *AdaptiveSelector {
	as := &AdaptiveSelector{
		weights:      weights,
		currentMood:  strings.ToLower(strings.TrimSpace(mood)),
		patterns:     patterns,
		currentHour:  hour,
		avoidanceMap: make(map[string]int),
	}

	if patterns != nil {
		// Build avoidance lookup map
		for _, entry := range patterns.AvoidanceList {
			as.avoidanceMap[entry.TaskText] = entry.TimesBypassed
		}

		// Resolve mood correlation
		if as.currentMood != "" {
			for _, mc := range patterns.MoodCorrelations {
				if mc.Mood == as.currentMood {
					as.preferredType = TaskType(mc.PreferredType)
					as.preferredEffort = TaskEffort(mc.PreferredEffort)
					break
				}
			}
		}

		// Find most productive period
		as.mostProductivePeriod = as.findMostProductivePeriod()
	}

	return as
}

// findMostProductivePeriod returns the period with the highest AvgTasksCompleted.
func (as *AdaptiveSelector) findMostProductivePeriod() string {
	if as.patterns == nil || len(as.patterns.TimeOfDayPatterns) == 0 {
		return ""
	}
	best := as.patterns.TimeOfDayPatterns[0]
	for _, p := range as.patterns.TimeOfDayPatterns[1:] {
		if p.AvgTasksCompleted > best.AvgTasksCompleted {
			best = p
		}
	}
	return best.Period
}

// ScoreCandidate computes a combined score for a candidate door set.
func (as *AdaptiveSelector) ScoreCandidate(candidate []*Task) float64 {
	diversity := float64(DiversityScore(candidate))
	mood := float64(MoodAlignmentScore(candidate, as.preferredType, as.preferredEffort))
	tod := as.timeOfDayBonus()
	avoidance := as.totalAvoidancePenalty(candidate)

	return diversity*as.weights.Diversity +
		mood*as.weights.Mood +
		tod*as.weights.TimeOfDay -
		avoidance*as.weights.Avoidance
}

// timeOfDayBonus returns a bonus if the current period is the most productive one.
func (as *AdaptiveSelector) timeOfDayBonus() float64 {
	if as.mostProductivePeriod == "" {
		return 0.0
	}
	currentPeriod := HourToPeriod(as.currentHour)
	if currentPeriod == as.mostProductivePeriod {
		return 1.0
	}
	return 0.0
}

// totalAvoidancePenalty sums the avoidance penalty for all tasks in the candidate set.
func (as *AdaptiveSelector) totalAvoidancePenalty(candidate []*Task) float64 {
	total := 0.0
	for _, t := range candidate {
		total += as.AvoidancePenalty(t)
	}
	return total
}

// AvoidancePenalty returns the penalty for a single task based on its bypass count.
func (as *AdaptiveSelector) AvoidancePenalty(task *Task) float64 {
	count, ok := as.avoidanceMap[task.Text]
	if !ok {
		return 0.0
	}
	if count >= 10 {
		return 0.8
	}
	if count >= 5 {
		return 0.5
	}
	return 0.0
}

// SelectDoorsAdaptive picks up to count tasks using adaptive scoring.
// Falls back to SelectDoors if selector is nil.
func SelectDoorsAdaptive(pool *TaskPool, count int, selector *AdaptiveSelector) []*Task {
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	return selectDoorsAdaptiveWithRand(pool, count, selector, rng)
}

// selectDoorsAdaptiveWithRand picks tasks using adaptive scoring with a deterministic RNG.
func selectDoorsAdaptiveWithRand(pool *TaskPool, count int, selector *AdaptiveSelector, rng *rand.Rand) []*Task {
	if selector == nil {
		return selectDoorsWithRand(pool, count, rng)
	}

	available := pool.GetAvailableForDoors()
	if len(available) == 0 {
		return nil
	}
	if len(available) <= count {
		for _, t := range available {
			pool.MarkRecentlyShown(t.ID)
		}
		return available
	}

	const numCandidates = 10

	bestScore := -1000.0 // allow negative scores from avoidance penalty
	var bestSet []*Task

	for i := range numCandidates {
		perm := make([]*Task, len(available))
		copy(perm, available)
		for j := range count {
			k := j + rng.IntN(len(perm)-j)
			perm[j], perm[k] = perm[k], perm[j]
		}
		candidate := perm[:count]
		score := selector.ScoreCandidate(candidate)

		if score > bestScore {
			bestScore = score
			bestSet = make([]*Task, count)
			copy(bestSet, candidate)
		} else if score == bestScore && rng.IntN(i+1) == 0 {
			bestSet = make([]*Task, count)
			copy(bestSet, candidate)
		}
	}

	// Diversity floor: if all tasks match same type, swap one out
	if selector.preferredType != "" {
		matchCount := 0
		for _, t := range bestSet {
			if t.Type == selector.preferredType {
				matchCount++
			}
		}
		if matchCount == count {
			bestSetIDs := make(map[string]bool, count)
			for _, t := range bestSet {
				bestSetIDs[t.ID] = true
			}
			var nonMatching []*Task
			for _, t := range available {
				if t.Type != selector.preferredType && !bestSetIDs[t.ID] {
					nonMatching = append(nonMatching, t)
				}
			}
			if len(nonMatching) > 0 {
				replacement := nonMatching[rng.IntN(len(nonMatching))]
				bestSet[count-1] = replacement
			}
		}
	}

	for _, t := range bestSet {
		pool.MarkRecentlyShown(t.ID)
	}
	return bestSet
}
