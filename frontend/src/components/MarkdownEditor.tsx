import { useState } from 'react';
import { MarkdownView } from './MarkdownView';
import { Code, Eye } from 'lucide-react';

interface Props {
  value: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  placeholder?: string;
  rows?: number;
}

export function MarkdownEditor({ value, onChange, onBlur, placeholder, rows = 8 }: Props) {
  const [preview, setPreview] = useState(false);

  const tabClass = (active: boolean) =>
    `flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md transition-colors ${
      active
        ? 'bg-surface-overlay text-text-primary'
        : 'text-text-muted hover:text-text-secondary'
    }`;

  return (
    <div className="border border-border-default rounded-lg overflow-hidden bg-surface-raised">
      <div className="flex items-center gap-1 px-2 py-1.5 border-b border-border-default bg-surface/50">
        <button type="button" onClick={() => setPreview(false)} className={tabClass(!preview)}>
          <Code className="w-3.5 h-3.5" />
          Редактор
        </button>
        <button type="button" onClick={() => setPreview(true)} className={tabClass(preview)}>
          <Eye className="w-3.5 h-3.5" />
          Превью
        </button>
      </div>

      {preview ? (
        <div className="p-4 min-h-[calc(1.5rem*8+1.25rem)]">
          {value ? (
            <MarkdownView content={value} />
          ) : (
            <p className="text-text-muted text-sm">Нет содержимого</p>
          )}
        </div>
      ) : (
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          rows={rows}
          className="w-full bg-transparent px-3 py-2.5 text-text-primary text-sm font-mono focus:outline-none resize-y"
          placeholder={placeholder}
          onBlur={onBlur}
        />
      )}
    </div>
  );
}
