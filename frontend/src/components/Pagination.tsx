import { ChevronLeft, ChevronRight } from 'lucide-react';

interface Props {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export function Pagination({ page, totalPages, onPageChange }: Props) {
  if (totalPages <= 1) return null;

  const btnClass =
    'flex items-center gap-1 px-3 py-1.5 rounded-lg text-sm bg-surface-raised border border-border-default text-text-secondary hover:border-border-hover hover:text-text-primary transition-colors disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:border-border-default';

  return (
    <div className="flex items-center gap-3 justify-center mt-8">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1}
        className={btnClass}
      >
        <ChevronLeft className="w-4 h-4" />
        Назад
      </button>
      <span className="text-sm text-text-muted font-mono tabular-nums">
        {page} / {totalPages}
      </span>
      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
        className={btnClass}
      >
        Вперёд
        <ChevronRight className="w-4 h-4" />
      </button>
    </div>
  );
}
