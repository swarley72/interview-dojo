import { useQuestionsStore } from '../stores/questions';
import { useTagsStore } from '../stores/tags';
import { Filter } from 'lucide-react';
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

export function FilterBar() {
  const { filterDifficulty, filterType, filterTagIds, setFilter, setFilterTagIds } =
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
