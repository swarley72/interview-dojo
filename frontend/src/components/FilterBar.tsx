import { useEffect, useRef, useState } from 'react';
import { useQuestionsStore } from '../stores/questions';
import { useTagsStore } from '../stores/tags';
import { Filter, Search, X } from 'lucide-react';
import { Dropdown, MultiDropdown } from './Dropdown';

const difficultyOptions = [
  { value: 'easy', label: 'Легко' },
  { value: 'medium', label: 'Средне' },
  { value: 'hard', label: 'Сложно' },
];

const typeOptions = [
  { value: 'theory', label: 'Theory' },
  { value: 'coding', label: 'Coding' },
  { value: 'algorithm', label: 'Algorithm' },
  { value: 'system_design', label: 'System Design' },
];

const verifiedOptions = [
  { value: 'true', label: 'Verified' },
  { value: 'false', label: 'Unverified' },
];

export function SearchInput() {
  const { searchQuery, setSearchQuery } = useQuestionsStore();
  const [localQuery, setLocalQuery] = useState(searchQuery);
  const timerRef = useRef<ReturnType<typeof setTimeout>>(null);

  useEffect(() => {
    setLocalQuery(searchQuery);
  }, [searchQuery]);

  const handleSearchChange = (value: string) => {
    setLocalQuery(value);
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => setSearchQuery(value), 300);
  };

  return (
    <div className="relative">
      <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-muted" />
      <input
        value={localQuery}
        onChange={(e) => handleSearchChange(e.target.value)}
        placeholder="Поиск..."
        className="pl-8 pr-8 py-1.5 text-sm bg-surface-raised border border-border-default rounded-lg text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors w-72"
      />
      {localQuery && (
        <button
          type="button"
          onClick={() => handleSearchChange('')}
          className="absolute right-2 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary transition-colors"
        >
          <X className="w-3.5 h-3.5" />
        </button>
      )}
    </div>
  );
}

export function FilterBar() {
  const { filterDifficulty, filterType, filterTagIds, filterVerified, setFilter, setFilterTagIds } =
    useQuestionsStore();
  const tags = useTagsStore((s) => s.tags);
  const tagOptions = tags.map((t) => ({ value: t.id, label: t.name }));

  return (
    <div className="flex items-center gap-3 flex-wrap">
      <Filter className="w-4 h-4 text-text-muted" />
      <Dropdown
        value={filterDifficulty}
        options={difficultyOptions}
        placeholder="Все сложности"
        onChange={(v) => setFilter('filterDifficulty', v)}
      />
      <Dropdown
        value={filterType}
        options={typeOptions}
        placeholder="Все типы"
        onChange={(v) => setFilter('filterType', v)}
      />
      <Dropdown
        value={filterVerified}
        options={verifiedOptions}
        placeholder="Все статусы"
        onChange={(v) => setFilter('filterVerified', v)}
      />
      {tags.length > 0 && (
        <MultiDropdown
          selected={filterTagIds}
          options={tagOptions}
          placeholder="Теги"
          onChange={setFilterTagIds}
        />
      )}
    </div>
  );
}
