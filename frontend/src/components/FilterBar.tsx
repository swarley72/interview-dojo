import { useQuestionsStore } from '../stores/questions';
import { Filter } from 'lucide-react';

export function FilterBar() {
  const { filterDifficulty, filterType, setFilter } = useQuestionsStore();

  const selectClass =
    'bg-surface-raised text-text-primary border border-border-default rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:border-accent transition-colors appearance-none cursor-pointer';

  return (
    <div className="flex items-center gap-3 flex-wrap">
      <Filter className="w-4 h-4 text-text-muted" />
      <select
        value={filterDifficulty}
        onChange={(e) => setFilter('filterDifficulty', e.target.value)}
        className={selectClass}
      >
        <option value="">Все сложности</option>
        <option value="easy">Easy</option>
        <option value="medium">Medium</option>
        <option value="hard">Hard</option>
      </select>
      <select
        value={filterType}
        onChange={(e) => setFilter('filterType', e.target.value)}
        className={selectClass}
      >
        <option value="">Все типы</option>
        <option value="theory">Theory</option>
        <option value="coding">Coding</option>
        <option value="algorithm">Algorithm</option>
        <option value="system_design">System Design</option>
      </select>
    </div>
  );
}
