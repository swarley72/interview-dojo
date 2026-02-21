import { Hash } from 'lucide-react';

export function TagBadge({ name }: { name: string }) {
  return (
    <span className="inline-flex items-center gap-0.5 px-2 py-0.5 rounded-md text-xs bg-surface-overlay text-text-secondary border border-border-default font-mono">
      <Hash className="w-3 h-3 text-text-muted" />
      {name}
    </span>
  );
}
