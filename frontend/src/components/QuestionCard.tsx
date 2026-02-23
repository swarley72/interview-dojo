import { Link } from 'react-router-dom';
import type { QuestionShort } from '../types';
import { DifficultyBadge } from './DifficultyBadge';
import { TypeBadge } from './TypeBadge';
import { TagBadge } from './TagBadge';
import { useTagsStore } from '../stores/tags';
import { ChevronRight, ShieldCheck } from 'lucide-react';

export function QuestionCard({ question }: { question: QuestionShort }) {
  const tagMap = useTagsStore((s) => s.tagMap);

  return (
    <Link
      to={`/questions/${question.id}`}
      className="group flex items-center justify-between bg-surface-raised rounded-xl p-4 border border-border-default hover:border-border-hover transition-all"
    >
      <div className="space-y-2 min-w-0">
        <h3 className="text-text-primary font-medium group-hover:text-accent transition-colors truncate">
          {question.title}
        </h3>
        <div className="flex gap-2 flex-wrap items-center">
          <DifficultyBadge difficulty={question.difficulty} />
          <TypeBadge type={question.type} />
          {question.verified && (
            <span className="flex items-center gap-1 text-xs text-green-400">
              <ShieldCheck className="w-3.5 h-3.5" />
            </span>
          )}
          {question.tag_ids?.map((id) => (
            <TagBadge key={id} name={tagMap[id] ?? `${id}`} />
          ))}
        </div>
      </div>
      <ChevronRight className="w-5 h-5 text-text-muted group-hover:text-accent transition-colors shrink-0 ml-4" />
    </Link>
  );
}
