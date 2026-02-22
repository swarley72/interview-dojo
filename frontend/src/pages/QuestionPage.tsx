import { useCallback, useEffect, useMemo, useState } from 'react';
import { useParams, useNavigate, useBlocker } from 'react-router-dom';
import type { Question, Difficulty, QuestionType } from '../types';
import { questionsApi } from '../api/questions';
import { useAuthStore } from '../stores/auth';
import { useTagsStore } from '../stores/tags';
import { QuestionView } from '../components/QuestionView';
import { MarkdownEditor } from '../components/MarkdownEditor';
import { ConfirmModal } from '../components/ConfirmModal';
import { TagSelector } from '../components/TagSelector';
import { Pencil, Save, X, Loader2, Trash2 } from 'lucide-react';

export function QuestionPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [question, setQuestion] = useState<Question | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const user = useAuthStore((s) => s.user);
  const fetchTags = useTagsStore((s) => s.fetchTags);

  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [editDifficulty, setEditDifficulty] = useState<Difficulty>('easy');
  const [editType, setEditType] = useState<QuestionType>('theory');
  const [editContent, setEditContent] = useState('');
  const [editAnswer, setEditAnswer] = useState('');
  const [editTagIds, setEditTagIds] = useState<number[]>([]);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const isDirty = useMemo(() => {
    if (!editing || !question) return false;
    return (
      editTitle !== question.title ||
      editDifficulty !== question.difficulty ||
      editType !== question.type ||
      editContent !== (question.content_md ?? '') ||
      editAnswer !== (question.answer_md ?? '') ||
      JSON.stringify(editTagIds) !== JSON.stringify(question.tag_ids ?? [])
    );
  }, [editing, question, editTitle, editDifficulty, editType, editContent, editAnswer, editTagIds]);

  const blocker = useBlocker(editing && isDirty);

  useEffect(() => {
    if (!editing || !isDirty) return;
    const handler = (e: BeforeUnloadEvent) => e.preventDefault();
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [editing, isDirty]);

  useEffect(() => {
    fetchTags();
    if (!id) return;
    setLoading(true);
    questionsApi
      .get(id)
      .then((q) => {
        setQuestion(q);
        setLoading(false);
      })
      .catch((e) => {
        setError((e as Error).message);
        setLoading(false);
      });
  }, [id]);

  const startEditing = () => {
    if (!question) return;
    setEditTitle(question.title);
    setEditDifficulty(question.difficulty);
    setEditType(question.type);
    setEditContent(question.content_md ?? '');
    setEditAnswer(question.answer_md ?? '');
    setEditTagIds(question.tag_ids ?? []);
    setEditing(true);
  };

  const [showUnsavedModal, setShowUnsavedModal] = useState(false);

  const cancelEditing = useCallback(() => {
    if (isDirty) {
      setShowUnsavedModal(true);
      return;
    }
    setEditing(false);
  }, [isDirty]);

  const saveEditing = async () => {
    if (!id) return;
    setSaving(true);
    try {
      const updated = await questionsApi.update(id, {
        title: editTitle,
        difficulty: editDifficulty,
        type: editType,
        content_md: editContent,
        answer_md: editAnswer,
        tag_ids: editTagIds,
      });
      setQuestion(updated);
      setEditing(false);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!id) return;
    setDeleting(true);
    try {
      await questionsApi.delete(id);
      navigate('/');
    } catch (e) {
      setError((e as Error).message);
      setShowDeleteModal(false);
    } finally {
      setDeleting(false);
    }
  };

  const inputClass =
    'w-full bg-surface-raised border border-border-default rounded-lg px-3 py-2.5 text-text-primary focus:outline-none focus:border-accent transition-colors';
  const selectClass =
    'bg-surface-raised border border-border-default rounded-lg px-3 py-1.5 text-sm text-text-primary focus:outline-none focus:border-accent transition-colors';

  if (loading)
    return (
      <div className="flex items-center gap-2 text-text-muted text-sm py-12 justify-center">
        <Loader2 className="w-4 h-4 animate-spin" />
        Загрузка...
      </div>
    );
  if (error) return <p className="text-red-400 text-sm">{error}</p>;
  if (!question) return <p className="text-text-muted text-sm">Вопрос не найден</p>;

  if (editing) {
    return (
      <div className="space-y-4 max-w-6xl">
        <h1 className="text-xl font-bold text-text-primary flex items-center gap-2">
          <Pencil className="w-5 h-5 text-accent" />
          Редактирование
        </h1>
        <input
          value={editTitle}
          onChange={(e) => setEditTitle(e.target.value)}
          className={inputClass}
          placeholder="Заголовок"
        />
        <div className="flex gap-3">
          <select
            value={editDifficulty}
            onChange={(e) => setEditDifficulty(e.target.value as Difficulty)}
            className={selectClass}
          >
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
          </select>
          <select
            value={editType}
            onChange={(e) => setEditType(e.target.value as QuestionType)}
            className={selectClass}
          >
            <option value="theory">Theory</option>
            <option value="coding">Coding</option>
            <option value="algorithm">Algorithm</option>
            <option value="system_design">System Design</option>
          </select>
        </div>
        <TagSelector selected={editTagIds} onChange={setEditTagIds} />
        <MarkdownEditor
          value={editContent}
          onChange={setEditContent}
          placeholder="Контент (Markdown)"
        />
        <MarkdownEditor
          value={editAnswer}
          onChange={setEditAnswer}
          placeholder="Ответ (Markdown)"
        />
        <div className="flex gap-3">
          <button
            onClick={saveEditing}
            disabled={saving}
            className="flex items-center gap-2 px-4 py-2 bg-accent hover:bg-accent-hover text-white rounded-lg text-sm font-semibold transition-colors disabled:opacity-50"
          >
            <Save className="w-4 h-4" />
            {saving ? 'Сохранение...' : 'Сохранить'}
          </button>
          <button
            onClick={cancelEditing}
            className="flex items-center gap-2 px-4 py-2 bg-surface-overlay hover:bg-border-hover text-text-secondary rounded-lg text-sm font-medium border border-border-default transition-colors"
          >
            <X className="w-4 h-4" />
            Отмена
          </button>
        </div>
        {error && <p className="text-red-400 text-sm">{error}</p>}

        <ConfirmModal
          open={showUnsavedModal}
          title="Несохранённые изменения"
          message="У вас есть несохранённые изменения. Уйти без сохранения?"
          confirmLabel="Уйти"
          onConfirm={() => {
            setShowUnsavedModal(false);
            setEditing(false);
          }}
          onCancel={() => setShowUnsavedModal(false)}
        />

        <ConfirmModal
          open={blocker.state === 'blocked'}
          title="Несохранённые изменения"
          message="У вас есть несохранённые изменения. Покинуть страницу без сохранения?"
          confirmLabel="Покинуть"
          onConfirm={() => blocker.proceed?.()}
          onCancel={() => blocker.reset?.()}
        />
      </div>
    );
  }

  return (
    <div className="max-w-6xl">
      {user?.is_super_user && (
        <div className="mb-4 flex items-center gap-2">
          <button
            onClick={startEditing}
            className="flex items-center gap-2 px-3 py-1.5 text-sm bg-surface-raised text-text-secondary border border-border-default hover:border-border-hover hover:text-text-primary rounded-lg transition-colors"
          >
            <Pencil className="w-3.5 h-3.5" />
            Редактировать
          </button>
          <button
            onClick={() => setShowDeleteModal(true)}
            className="flex items-center gap-2 px-3 py-1.5 text-sm bg-surface-raised text-red-400 border border-border-default hover:border-red-500/30 hover:bg-red-500/10 rounded-lg transition-colors"
          >
            <Trash2 className="w-3.5 h-3.5" />
            Удалить
          </button>
        </div>
      )}
      <QuestionView question={question} />

      <ConfirmModal
        open={showDeleteModal}
        title="Удалить вопрос"
        message="Вопрос будет удалён без возможности восстановления. Продолжить?"
        confirmLabel="Удалить"
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteModal(false)}
        loading={deleting}
      />
    </div>
  );
}
