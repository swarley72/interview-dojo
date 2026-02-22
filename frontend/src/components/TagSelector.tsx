import { useState } from 'react';
import { useTagsStore } from '../stores/tags';
import { tagsApi } from '../api/tags';
import { ApiError } from '../api/client';
import { ConfirmModal } from './ConfirmModal';
import { Plus, X } from 'lucide-react';

interface Props {
  selected: number[];
  onChange: (ids: number[]) => void;
}

export function TagSelector({ selected, onChange }: Props) {
  const { tags, fetchTags } = useTagsStore();
  const [newTag, setNewTag] = useState('');
  const [creating, setCreating] = useState(false);
  const [deleteTag, setDeleteTag] = useState<{ id: number; name: string } | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState('');

  const toggle = (tagId: number) => {
    onChange(
      selected.includes(tagId)
        ? selected.filter((id) => id !== tagId)
        : [...selected, tagId],
    );
  };

  const handleCreate = async () => {
    const name = newTag.trim();
    if (!name) return;
    setCreating(true);
    setError('');
    try {
      const tag = await tagsApi.create(name);
      await fetchTags();
      onChange([...selected, tag.id]);
      setNewTag('');
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setError('Тег уже существует');
      } else {
        setError('Ошибка создания тега');
      }
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async () => {
    if (!deleteTag) return;
    setDeleting(true);
    try {
      await tagsApi.delete(deleteTag.id);
      onChange(selected.filter((id) => id !== deleteTag.id));
      await fetchTags();
      setDeleteTag(null);
    } catch {
      // ignore
    } finally {
      setDeleting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleCreate();
    }
  };

  return (
    <>
      <div className="flex gap-2 flex-wrap items-center">
        {tags.map((tag) => (
          <div
            key={tag.id}
            className={`group inline-flex items-center gap-1 rounded-lg text-xs font-medium border transition-colors ${
              selected.includes(tag.id)
                ? 'bg-accent/10 text-accent border-accent/30'
                : 'bg-surface-overlay text-text-muted border-border-default hover:border-border-hover'
            }`}
          >
            <button
              type="button"
              onClick={() => toggle(tag.id)}
              className="pl-2.5 py-1"
            >
              {tag.name}
            </button>
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                setDeleteTag({ id: tag.id, name: tag.name });
              }}
              className="pr-1.5 py-1 opacity-0 group-hover:opacity-100 text-red-400 hover:text-red-500 transition-all"
            >
              <X className="w-3 h-3" />
            </button>
          </div>
        ))}
        <div className="flex items-center gap-1">
          <input
            type="text"
            value={newTag}
            onChange={(e) => setNewTag(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Новый тег..."
            className="w-28 bg-surface-overlay border border-border-default rounded-lg px-2.5 py-1 text-xs text-text-primary placeholder-text-muted focus:outline-none focus:border-accent transition-colors"
          />
          <button
            type="button"
            onClick={handleCreate}
            disabled={creating || !newTag.trim()}
            className="p-1 rounded-lg bg-accent/10 text-accent hover:bg-accent/20 transition-colors disabled:opacity-40"
          >
            <Plus className="w-3.5 h-3.5" />
          </button>
        </div>
        {error && (
          <span className="text-xs text-red-400">{error}</span>
        )}
      </div>

      <ConfirmModal
        open={!!deleteTag}
        title="Удалить тег"
        message={`Тег «${deleteTag?.name}» будет удалён у всех вопросов. Продолжить?`}
        confirmLabel="Удалить"
        onConfirm={handleDelete}
        onCancel={() => setDeleteTag(null)}
        loading={deleting}
      />
    </>
  );
}
