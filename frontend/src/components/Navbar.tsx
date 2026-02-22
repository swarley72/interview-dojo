import { useState } from 'react';
import { Link, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../stores/auth';
import { ankiApi } from '../api/anki';
import { ConfirmModal } from './ConfirmModal';
import { BookOpen, RefreshCw, LogOut, Flame, RotateCcw } from 'lucide-react';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();
  const [showReset, setShowReset] = useState(false);
  const [resetting, setResetting] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleReset = async () => {
    setResetting(true);
    try {
      await ankiApi.resetProgress();
      setShowReset(false);
      window.location.reload();
    } catch {
      // ignore
    } finally {
      setResetting(false);
    }
  };

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `flex items-center gap-1.5 text-sm transition-colors ${
      isActive
        ? 'text-accent font-semibold'
        : 'text-text-secondary hover:text-text-primary'
    }`;

  return (
    <>
      <nav className="border-b border-border-default bg-surface/80 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
          <div className="flex items-center gap-8">
            <Link to="/" className="flex items-center gap-2 group">
              <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center group-hover:bg-accent/20 transition-colors">
                <Flame className="w-4 h-4 text-accent" />
              </div>
              <span className="text-base font-bold text-text-primary tracking-tight">
                Interview<span className="text-accent">Dojo</span>
              </span>
            </Link>
            {isAuthenticated && (
              <div className="flex items-center gap-6">
                <NavLink to="/" className={linkClass} end>
                  <BookOpen className="w-4 h-4" />
                  Каталог
                </NavLink>
                <NavLink to="/review" className={linkClass}>
                  <RefreshCw className="w-4 h-4" />
                  Повторение
                </NavLink>
              </div>
            )}
          </div>
          {isAuthenticated && user && (
            <div className="flex items-center gap-4">
              <button
                onClick={() => setShowReset(true)}
                className="flex items-center gap-1.5 text-sm text-text-muted hover:text-red-400 transition-colors"
                title="Сбросить прогресс повторения"
              >
                <RotateCcw className="w-3.5 h-3.5" />
                Сбросить
              </button>
              <span className="text-sm text-text-muted font-mono">{user.login}</span>
              <button
                onClick={handleLogout}
                className="flex items-center gap-1.5 text-sm text-text-muted hover:text-red-400 transition-colors"
              >
                <LogOut className="w-4 h-4" />
                Выйти
              </button>
            </div>
          )}
        </div>
      </nav>
      <ConfirmModal
        open={showReset}
        title="Сбросить прогресс"
        message="Весь прогресс повторения будет удалён. Все вопросы начнутся с нуля. Это действие необратимо."
        confirmLabel="Сбросить"
        loading={resetting}
        onConfirm={handleReset}
        onCancel={() => setShowReset(false)}
      />
    </>
  );
}
