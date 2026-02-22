import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Difficulty, QuestionType } from '../types';
import { questionsApi } from '../api/questions';
import { useTagsStore } from '../stores/tags';
import { MarkdownEditor } from '../components/MarkdownEditor';
import { ExcalidrawEditor } from '../components/ExcalidrawEditor';
import { TagSelector } from '../components/TagSelector';
import { Plus, Save } from 'lucide-react';

export function CreateQuestionPage() {
  const navigate = useNavigate();
  const fetchTags = useTagsStore((s) => s.fetchTags);

  const [title, setTitle] = useState('');
  const [difficulty, setDifficulty] = useState<Difficulty>('easy');
  const [type, setType] = useState<QuestionType>('theory');
  const [contentMd, setContentMd] = useState('');
  const [answerMd, setAnswerMd] = useState('');
  const [tagIds, setTagIds] = useState<number[]>([]);
  const [excalidrawJson, setExcalidrawJson] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchTags();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return setError('Заголовок обязателен');
    setSaving(true);
    setError('');
    try {
      const q = await questionsApi.create({
        title,
        difficulty,
        type,
        content_md: contentMd || undefined,
        answer_md: answerMd || undefined,
        excalidraw_json: excalidrawJson || undefined,
        tag_ids: tagIds.length > 0 ? tagIds : undefined,
      });
      navigate(`/questions/${q.id}`);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  const inputClass =
    'w-full bg-surface-raised border border-border-default rounded-lg px-3 py-2.5 text-text-primary focus:outline-none focus:border-accent transition-colors';
  const selectClass =
    'bg-surface-raised border border-border-default rounded-lg px-3 py-1.5 text-sm text-text-primary focus:outline-none focus:border-accent transition-colors';

  return (
    <div className="max-w-6xl">
      <h1 className="text-xl font-bold text-text-primary mb-6 flex items-center gap-2">
        <Plus className="w-5 h-5 text-accent" />
        Новый вопрос
      </h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-sm rounded-lg p-3">
            {error}
          </div>
        )}

        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className={inputClass}
          placeholder="Заголовок"
        />

        <div className="flex gap-3">
          <select
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value as Difficulty)}
            className={selectClass}
          >
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
          </select>
          <select
            value={type}
            onChange={(e) => setType(e.target.value as QuestionType)}
            className={selectClass}
          >
            <option value="theory">Theory</option>
            <option value="coding">Coding</option>
            <option value="algorithm">Algorithm</option>
            <option value="system_design">System Design</option>
          </select>
        </div>

        <TagSelector selected={tagIds} onChange={setTagIds} />

        <MarkdownEditor
          value={contentMd}
          onChange={setContentMd}
          placeholder="Контент (Markdown)"
        />

        <MarkdownEditor
          value={answerMd}
          onChange={setAnswerMd}
          placeholder="Ответ (Markdown)"
        />

        {type === 'system_design' && (
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-2">Диаграмма</label>
            <ExcalidrawEditor initialData={excalidrawJson || null} onChange={setExcalidrawJson} />
          </div>
        )}

        <button
          type="submit"
          disabled={saving}
          className="flex items-center gap-2 px-4 py-2.5 bg-accent hover:bg-accent-hover text-white rounded-lg text-sm font-semibold transition-colors disabled:opacity-50"
        >
          <Save className="w-4 h-4" />
          {saving ? 'Создание...' : 'Создать'}
        </button>
      </form>
    </div>
  );
}
