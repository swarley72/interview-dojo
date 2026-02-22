import type { Difficulty } from '../types';
import { Signal, SignalMedium, SignalLow } from 'lucide-react';

const config: Record<Difficulty, { class: string; icon: typeof Signal; label: string }> = {
  easy: { class: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20', icon: SignalLow, label: 'Легко' },
  medium: { class: 'bg-amber-500/10 text-amber-400 border-amber-500/20', icon: SignalMedium, label: 'Средне' },
  hard: { class: 'bg-red-500/10 text-red-400 border-red-500/20', icon: Signal, label: 'Сложно' },
};

export function DifficultyBadge({ difficulty }: { difficulty: Difficulty }) {
  const c = config[difficulty];
  const Icon = c.icon;
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-medium border ${c.class}`}>
      <Icon className="w-3 h-3" />
      {c.label}
    </span>
  );
}
