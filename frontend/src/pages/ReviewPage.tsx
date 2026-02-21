import { useEffect, useState } from 'react';
import type { AnswerQuality, Question, UserProgressResponse } from '../types';
import { ankiApi } from '../api/anki';
import { QuestionView } from '../components/QuestionView';
import { AnswerRating } from '../components/AnswerRating';
import { ApiError } from '../api/client';
import { RefreshCw, Loader2, PartyPopper, BarChart3 } from 'lucide-react';

export function ReviewPage() {
  const [question, setQuestion] = useState<Question | null>(null);
  const [loading, setLoading] = useState(true);
  const [empty, setEmpty] = useState(false);
  const [result, setResult] = useState<UserProgressResponse | null>(null);
  const [rating, setRating] = useState(false);

  const loadNext = async () => {
    setLoading(true);
    setResult(null);
    setEmpty(false);
    try {
      const q = await ankiApi.getNext();
      setQuestion(q);
    } catch (e) {
      if (e instanceof ApiError && e.status === 404) {
        setEmpty(true);
        setQuestion(null);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadNext();
  }, []);

  const handleRate = async (quality: AnswerQuality) => {
    if (!question) return;
    setRating(true);
    try {
      const res = await ankiApi.recordAnswer(question.id, quality);
      setResult(res);
      setTimeout(loadNext, 1500);
    } catch {
      // ignore
    } finally {
      setRating(false);
    }
  };

  if (loading)
    return (
      <div className="flex items-center gap-2 text-text-muted text-sm py-12 justify-center">
        <Loader2 className="w-4 h-4 animate-spin" />
        Загрузка...
      </div>
    );

  if (empty)
    return (
      <div className="text-center py-16">
        <PartyPopper className="w-12 h-12 text-accent mx-auto mb-4" />
        <p className="text-text-primary text-lg font-semibold">Нет вопросов для повторения</p>
        <p className="text-text-muted text-sm mt-1">Возвращайтесь позже!</p>
      </div>
    );

  if (!question) return null;

  return (
    <div className="max-w-6xl">
      <h1 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
        <RefreshCw className="w-5 h-5 text-accent" />
        Повторение
      </h1>
      <QuestionView question={question}>
        <div className="mt-6 space-y-4">
          <AnswerRating onRate={handleRate} disabled={rating || !!result} />
          {result && (
            <div className="bg-surface-raised rounded-xl p-4 border border-border-default text-sm flex items-center gap-3">
              <BarChart3 className="w-4 h-4 text-accent shrink-0" />
              <span className="text-text-secondary">
                {result.repetitions} повт., интервал {result.interval_days} дн., ease{' '}
                {result.ease_factor.toFixed(2)}
              </span>
            </div>
          )}
        </div>
      </QuestionView>
    </div>
  );
}
