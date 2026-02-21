import { useState } from 'react';
import type { Question } from '../types';
import { DifficultyBadge } from './DifficultyBadge';
import { TypeBadge } from './TypeBadge';
import { TagBadge } from './TagBadge';
import { MarkdownView } from './MarkdownView';
import { useTagsStore } from '../stores/tags';
import { Eye, BarChart3, Repeat, Calendar, Gauge, Timer } from 'lucide-react';

interface Props {
  question: Question;
  children?: React.ReactNode;
}

export function QuestionView({ question, children }: Props) {
  const [showAnswer, setShowAnswer] = useState(false);
  const tagMap = useTagsStore((s) => s.tagMap);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-primary mb-3">{question.title}</h1>
        <div className="flex gap-2 flex-wrap">
          <DifficultyBadge difficulty={question.difficulty} />
          <TypeBadge type={question.type} />
          {question.tag_ids?.map((id) => (
            <TagBadge key={id} name={tagMap[id] ?? `${id}`} />
          ))}
        </div>
      </div>

      {question.content_md && (
        <div className="bg-surface-raised rounded-xl p-6 border border-border-default">
          <MarkdownView content={question.content_md} />
        </div>
      )}

      {question.answer_md && (
        <div>
          {!showAnswer ? (
            <button
              onClick={() => setShowAnswer(true)}
              className="flex items-center gap-2 px-4 py-2.5 bg-accent hover:bg-accent-hover text-white rounded-lg text-sm font-semibold transition-colors"
            >
              <Eye className="w-4 h-4" />
              Показать ответ
            </button>
          ) : (
            <div className="bg-surface-raised rounded-xl p-6 border border-border-default">
              <h2 className="text-lg font-semibold text-text-primary mb-3 flex items-center gap-2">
                <Eye className="w-5 h-5 text-accent" />
                Ответ
              </h2>
              <MarkdownView content={question.answer_md} />
            </div>
          )}
        </div>
      )}

      {question.progress ? (
        <div className="bg-surface-raised rounded-xl p-4 border border-border-default">
          <h3 className="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
            <BarChart3 className="w-4 h-4 text-accent" />
            Прогресс
          </h3>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-sm">
            <div className="flex items-center gap-2">
              <Repeat className="w-4 h-4 text-text-muted" />
              <div>
                <div className="text-text-muted text-xs">Повторений</div>
                <div className="text-text-primary font-mono">{question.progress.repetitions}</div>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Timer className="w-4 h-4 text-text-muted" />
              <div>
                <div className="text-text-muted text-xs">Интервал</div>
                <div className="text-text-primary font-mono">{question.progress.interval_days} дн.</div>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Gauge className="w-4 h-4 text-text-muted" />
              <div>
                <div className="text-text-muted text-xs">Ease</div>
                <div className="text-text-primary font-mono">{question.progress.ease_factor.toFixed(2)}</div>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4 text-text-muted" />
              <div>
                <div className="text-text-muted text-xs">Следующее</div>
                <div className="text-text-primary font-mono">
                  {new Date(question.progress.next_review_at).toLocaleDateString('ru-RU')}
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <p className="text-sm text-text-muted">Ещё не отвечали</p>
      )}

      {children}
    </div>
  );
}
