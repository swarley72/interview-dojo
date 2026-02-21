import type { AnswerQuality } from '../types';
import { RotateCcw, Frown, Smile, Zap } from 'lucide-react';

interface Props {
  onRate: (quality: AnswerQuality) => void;
  disabled?: boolean;
}

const buttons: { quality: AnswerQuality; label: string; icon: typeof Zap; style: string }[] = [
  { quality: 'again', label: 'Снова', icon: RotateCcw, style: 'bg-red-500/10 text-red-400 border-red-500/20 hover:bg-red-500/20' },
  { quality: 'hard', label: 'Сложно', icon: Frown, style: 'bg-amber-500/10 text-amber-400 border-amber-500/20 hover:bg-amber-500/20' },
  { quality: 'good', label: 'Хорошо', icon: Smile, style: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 hover:bg-emerald-500/20' },
  { quality: 'easy', label: 'Легко', icon: Zap, style: 'bg-blue-500/10 text-blue-400 border-blue-500/20 hover:bg-blue-500/20' },
];

export function AnswerRating({ onRate, disabled }: Props) {
  return (
    <div className="flex gap-3 flex-wrap">
      {buttons.map((b) => {
        const Icon = b.icon;
        return (
          <button
            key={b.quality}
            onClick={() => onRate(b.quality)}
            disabled={disabled}
            className={`flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-semibold border transition-colors disabled:opacity-40 ${b.style}`}
          >
            <Icon className="w-4 h-4" />
            {b.label}
          </button>
        );
      })}
    </div>
  );
}
