import type { QuestionType } from '../types';
import { BookOpen, Code, Cpu, Network } from 'lucide-react';

const config: Record<QuestionType, { label: string; icon: typeof BookOpen }> = {
  theory: { label: 'theory', icon: BookOpen },
  coding: { label: 'coding', icon: Code },
  algorithm: { label: 'algorithm', icon: Cpu },
  system_design: { label: 'system design', icon: Network },
};

export function TypeBadge({ type }: { type: QuestionType }) {
  const c = config[type];
  const Icon = c.icon;
  return (
    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-medium bg-violet-500/10 text-violet-400 border border-violet-500/20">
      <Icon className="w-3 h-3" />
      {c.label}
    </span>
  );
}
