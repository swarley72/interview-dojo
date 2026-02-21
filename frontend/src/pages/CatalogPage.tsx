import { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useQuestionsStore } from '../stores/questions';
import { useTagsStore } from '../stores/tags';
import { useAuthStore } from '../stores/auth';
import { FilterBar } from '../components/FilterBar';
import { QuestionCard } from '../components/QuestionCard';
import { Pagination } from '../components/Pagination';
import { BookOpen, Loader2, SearchX, Plus } from 'lucide-react';

export function CatalogPage() {
  const { items, totalCount, page, limit, isLoading, fetchQuestions, setPage } =
    useQuestionsStore();
  const fetchTags = useTagsStore((s) => s.fetchTags);
  const user = useAuthStore((s) => s.user);

  useEffect(() => {
    fetchTags();
    fetchQuestions();
  }, []);

  const totalPages = Math.ceil(totalCount / limit);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <h1 className="text-xl font-bold text-text-primary flex items-center gap-2">
            <BookOpen className="w-5 h-5 text-accent" />
            Каталог вопросов
          </h1>
          {user?.is_super_user && (
            <Link
              to="/questions/new"
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-accent hover:bg-accent-hover text-white rounded-lg font-semibold transition-colors"
            >
              <Plus className="w-4 h-4" />
              Новый
            </Link>
          )}
        </div>
        <FilterBar />
      </div>

      {isLoading ? (
        <div className="flex items-center gap-2 text-text-muted text-sm py-12 justify-center">
          <Loader2 className="w-4 h-4 animate-spin" />
          Загрузка...
        </div>
      ) : items.length === 0 ? (
        <div className="text-center py-12">
          <SearchX className="w-10 h-10 text-text-muted mx-auto mb-3" />
          <p className="text-text-muted text-sm">Вопросов не найдено</p>
        </div>
      ) : (
        <div className="grid gap-3">
          {items.map((q) => (
            <QuestionCard key={q.id} question={q} />
          ))}
        </div>
      )}

      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
    </div>
  );
}
