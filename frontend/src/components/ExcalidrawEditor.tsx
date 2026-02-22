import { Excalidraw } from '@excalidraw/excalidraw';
import '@excalidraw/excalidraw/index.css';
import type { ExcalidrawImperativeAPI } from '@excalidraw/excalidraw/types';
import { useRef, useCallback, useState, useEffect } from 'react';
import { Maximize2, Minimize2 } from 'lucide-react';
import libraryData from '../assets/software-architecture.json';

interface Props {
  initialData?: string | null;
  onChange?: (data: string) => void;
  readOnly?: boolean;
}

export function ExcalidrawEditor({ initialData, onChange, readOnly = false }: Props) {
  const apiRef = useRef<ExcalidrawImperativeAPI | null>(null);
  const [fullscreen, setFullscreen] = useState(false);

  const parsed = initialData ? JSON.parse(initialData) : undefined;

  useEffect(() => {
    if (fullscreen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => { document.body.style.overflow = ''; };
  }, [fullscreen]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && fullscreen) setFullscreen(false);
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [fullscreen]);

  const handleChange = useCallback(() => {
    if (readOnly || !apiRef.current || !onChange) return;
    const elements = apiRef.current.getSceneElements();
    const state = apiRef.current.getAppState();
    onChange(
      JSON.stringify({
        elements,
        appState: {
          viewBackgroundColor: state.viewBackgroundColor,
          gridSize: state.gridSize,
        },
      }),
    );
  }, [readOnly, onChange]);

  return (
    <div
      className={
        fullscreen
          ? 'fixed inset-0 z-[9999] bg-white'
          : 'relative border border-border-default rounded-xl overflow-hidden'
      }
      style={fullscreen ? undefined : { height: 500 }}
    >
      <button
        type="button"
        onClick={() => setFullscreen(!fullscreen)}
        className="absolute top-2 right-2 z-[10000] p-1.5 rounded-md bg-surface-overlay border border-border-default text-text-muted hover:text-text-primary transition-colors"
        title={fullscreen ? 'Свернуть (Esc)' : 'На весь экран'}
      >
        {fullscreen ? <Minimize2 className="w-4 h-4" /> : <Maximize2 className="w-4 h-4" />}
      </button>
      <Excalidraw
        excalidrawAPI={(api) => { apiRef.current = api; }}
        initialData={{
          elements: parsed?.elements,
          appState: parsed?.appState,
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          libraryItems: readOnly ? undefined : (libraryData.library as any),
        }}
        onChange={handleChange}
        viewModeEnabled={readOnly}
        UIOptions={{
          canvasActions: {
            changeViewBackgroundColor: !readOnly,
            loadScene: false,
            saveToActiveFile: false,
            export: false,
          },
          tools: { image: false },
        }}
        libraryReturnUrl={readOnly ? undefined : window.location.href}
      />
    </div>
  );
}
