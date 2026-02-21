import { Link } from 'react-router-dom';
import { Ghost, ArrowLeft } from 'lucide-react';

export function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-surface">
      <div className="text-center">
        <Ghost className="w-16 h-16 text-text-muted mx-auto mb-4" />
        <h1 className="text-5xl font-bold text-text-muted font-mono">404</h1>
        <p className="text-text-secondary mt-3">Страница не найдена</p>
        <Link
          to="/"
          className="inline-flex items-center gap-2 mt-6 px-4 py-2.5 bg-accent hover:bg-accent-hover text-white rounded-lg text-sm font-semibold transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          На главную
        </Link>
      </div>
    </div>
  );
}
