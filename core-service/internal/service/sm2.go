package service

import "math"

type sm2Result struct {
	Repetitions  int32
	EaseFactor   float64
	IntervalDays int32
}

func calculateSM2(repetitions int32, easeFactor float64, intervalDays int32, quality int) sm2Result {
	if easeFactor == 0 {
		easeFactor = 2.5
	}

	if quality < 3 {
		return sm2Result{
			Repetitions:  0,
			EaseFactor:   easeFactor,
			IntervalDays: 0,
		}
	}

	newEF := easeFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	newEF = math.Max(newEF, 1.3)

	newReps := repetitions + 1

	var newInterval int32
	switch newReps {
	case 1:
		newInterval = 1
	case 2:
		newInterval = 6
	default:
		newInterval = int32(math.Round(float64(intervalDays) * newEF))
	}

	return sm2Result{
		Repetitions:  newReps,
		EaseFactor:   newEF,
		IntervalDays: newInterval,
	}
}
