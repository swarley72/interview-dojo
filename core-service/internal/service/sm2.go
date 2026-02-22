package service

import "math"

type sm2Result struct {
	Repetitions  int32
	EaseFactor   float64
	IntervalDays int32
}

// calculateSM2 реализует алгоритм SuperMemo 2 для расчёта интервалов повторения.
//
// Параметры:
//   - repetitions  — сколько раз подряд ответ был успешным (quality >= 3)
//   - easeFactor   — множитель интервала, начинается с 2.0, максимум 2.0
//   - intervalDays — текущий интервал в днях до следующего повторения
//   - quality      — оценка ответа: 0 (AGAIN), 3 (HARD), 4 (GOOD), 5 (EASY)
//
// Примеры цепочек повторений (все ответы GOOD, quality=4, EF=2.0):
//
//	Повтор 1 → через 1 день    (фиксированный)
//	Повтор 2 → через 3 дня     (фиксированный)
//	Повтор 3 → через 6 дней    (3 × 2.0)
//	Повтор 4 → через 12 дней   (6 × 2.0)
//	Повтор 5 → через 24 дня    (12 × 2.0)
//
// Если ответ AGAIN (quality=0) — сброс: repetitions=0, interval=0, вопрос
// появится при ближайшем review. EaseFactor не меняется.
//
// Если ответ HARD (quality=3) — EF снижается, интервалы растут медленнее.
// Если ответ EASY (quality=5) — EF повышается, но не выше 2.0.
func calculateSM2(repetitions int32, easeFactor float64, intervalDays int32, quality int) sm2Result {
	if easeFactor == 0 {
		easeFactor = 2.0
	}

	// Плохой ответ — полный сброс прогресса, вопрос начинается с нуля.
	// EF сохраняется, чтобы не терять накопленную "сложность" карточки.
	if quality < 3 {
		return sm2Result{
			Repetitions:  0,
			EaseFactor:   easeFactor,
			IntervalDays: 0,
		}
	}

	// Пересчёт EaseFactor по формуле SM-2:
	//   quality=5 (EASY) → EF += 0.10  (интервалы растут быстрее)
	//   quality=4 (GOOD) → EF += 0.00  (без изменений)
	//   quality=3 (HARD) → EF -= 0.14  (интервалы растут медленнее)
	newEF := easeFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	newEF = math.Min(newEF, 2.0)

	// repetitions — счётчик успешных ответов подряд.
	// Определяет стадию повторения и формулу для интервала.
	newReps := repetitions + 1

	var newInterval int32
	switch newReps {
	case 1:
		// Первый успешный ответ → повторить завтра.
		newInterval = 1
	case 2:
		// Второй подряд → повторить через 3 дня.
		newInterval = 3
	default:
		// Далее: новый интервал = предыдущий × EaseFactor.
		// Пример: intervalDays=3, EF=2.0 → 3×2.0 = 6 дней.
		newInterval = int32(math.Round(float64(intervalDays) * newEF))
	}

	return sm2Result{
		Repetitions:  newReps,
		EaseFactor:   newEF,
		IntervalDays: newInterval,
	}
}
