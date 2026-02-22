import { useState, useRef, useEffect } from 'react';
import { ChevronDown, X } from 'lucide-react';

interface Option {
  value: string;
  label: string;
}

interface DropdownProps {
  value: string;
  options: Option[];
  placeholder: string;
  onChange: (value: string) => void;
}

export function Dropdown({ value, options, placeholder, onChange }: DropdownProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  const selected = options.find((o) => o.value === value);

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1.5 bg-surface-raised border border-border-default rounded-lg px-3 py-1.5 text-sm transition-colors hover:border-border-hover focus:outline-none focus:border-accent cursor-pointer min-w-[140px]"
      >
        <span className={selected ? 'text-text-primary' : 'text-text-muted'}>
          {selected ? selected.label : placeholder}
        </span>
        <span className="ml-auto flex items-center gap-1">
          {value && (
            <span
              role="button"
              onClick={(e) => {
                e.stopPropagation();
                onChange('');
                setOpen(false);
              }}
              className="text-text-muted hover:text-text-secondary transition-colors"
            >
              <X className="w-3.5 h-3.5" />
            </span>
          )}
          <ChevronDown
            className={`w-3.5 h-3.5 text-text-muted transition-transform ${open ? 'rotate-180' : ''}`}
          />
        </span>
      </button>

      {open && (
        <div className="absolute z-50 mt-1 w-full min-w-[160px] bg-surface-raised border border-border-default rounded-lg shadow-lg py-1 animate-in">
          {options.map((opt) => (
            <button
              key={opt.value}
              type="button"
              onClick={() => {
                onChange(opt.value);
                setOpen(false);
              }}
              className={`w-full text-left px-3 py-1.5 text-sm transition-colors cursor-pointer ${
                opt.value === value
                  ? 'bg-accent/10 text-accent font-medium'
                  : 'text-text-primary hover:bg-surface-overlay'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

interface MultiDropdownProps {
  selected: number[];
  options: { value: number; label: string }[];
  placeholder: string;
  onChange: (values: number[]) => void;
}

export function MultiDropdown({ selected, options, placeholder, onChange }: MultiDropdownProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const ref = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
        setSearch('');
      }
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  useEffect(() => {
    if (open && inputRef.current) inputRef.current.focus();
  }, [open]);

  const filtered = options.filter((o) =>
    o.label.toLowerCase().includes(search.toLowerCase()),
  );

  const toggle = (val: number) => {
    onChange(
      selected.includes(val) ? selected.filter((v) => v !== val) : [...selected, val],
    );
  };

  const selectedLabels = options.filter((o) => selected.includes(o.value));

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1.5 bg-surface-raised border border-border-default rounded-lg px-3 py-1.5 text-sm transition-colors hover:border-border-hover focus:outline-none focus:border-accent cursor-pointer min-w-[140px]"
      >
        {selectedLabels.length > 0 ? (
          <span className="flex items-center gap-1 flex-wrap">
            <span className="text-text-primary">{placeholder}</span>
            <span className="bg-accent/10 text-accent text-xs font-medium px-1.5 py-0.5 rounded">
              {selectedLabels.length}
            </span>
          </span>
        ) : (
          <span className="text-text-muted">{placeholder}</span>
        )}
        <span className="ml-auto flex items-center gap-1">
          {selected.length > 0 && (
            <span
              role="button"
              onClick={(e) => {
                e.stopPropagation();
                onChange([]);
                setOpen(false);
              }}
              className="text-text-muted hover:text-text-secondary transition-colors"
            >
              <X className="w-3.5 h-3.5" />
            </span>
          )}
          <ChevronDown
            className={`w-3.5 h-3.5 text-text-muted transition-transform ${open ? 'rotate-180' : ''}`}
          />
        </span>
      </button>

      {open && (
        <div className="absolute z-50 mt-1 w-full min-w-[200px] bg-surface-raised border border-border-default rounded-lg shadow-lg animate-in">
          <div className="p-2 border-b border-border-default">
            <input
              ref={inputRef}
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Поиск..."
              className="w-full bg-surface-overlay border border-border-default rounded px-2 py-1 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
            />
          </div>
          <div className="max-h-[200px] overflow-y-auto py-1">
            {filtered.length === 0 ? (
              <div className="px-3 py-2 text-sm text-text-muted">Не найдено</div>
            ) : (
              filtered.map((opt) => {
                const isSelected = selected.includes(opt.value);
                return (
                  <button
                    key={opt.value}
                    type="button"
                    onClick={() => toggle(opt.value)}
                    className={`w-full text-left px-3 py-1.5 text-sm transition-colors cursor-pointer flex items-center gap-2 ${
                      isSelected
                        ? 'bg-accent/10 text-accent'
                        : 'text-text-primary hover:bg-surface-overlay'
                    }`}
                  >
                    <span
                      className={`w-3.5 h-3.5 rounded border flex-shrink-0 flex items-center justify-center ${
                        isSelected
                          ? 'bg-accent border-accent'
                          : 'border-border-default'
                      }`}
                    >
                      {isSelected && (
                        <svg viewBox="0 0 12 12" className="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" strokeWidth="2">
                          <path d="M2 6l3 3 5-5" />
                        </svg>
                      )}
                    </span>
                    {opt.label}
                  </button>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
}
