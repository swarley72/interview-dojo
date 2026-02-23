import { useCallback, useRef, useState } from 'react';

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';

interface UseAutosaveOptions<T> {
  id: string | undefined;
  editing: boolean;
  isDirty: boolean;
  getPayload: () => T;
  onSaved: (result: unknown) => void;
  onError: (err: Error) => void;
  saveFn: (id: string, payload: T) => Promise<unknown>;
}

export function useAutosave<T>({
  id,
  editing,
  isDirty,
  getPayload,
  onSaved,
  onError,
  saveFn,
}: UseAutosaveOptions<T>) {
  const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle');
  const seqRef = useRef(0);
  const timerRef = useRef<ReturnType<typeof setTimeout>>(null);
  const isDirtyRef = useRef(isDirty);
  isDirtyRef.current = isDirty;
  const editingRef = useRef(editing);
  editingRef.current = editing;

  const triggerSave = useCallback(() => {
    if (!id || !editingRef.current || !isDirtyRef.current) return;
    const seq = ++seqRef.current;
    setSaveStatus('saving');

    saveFn(id, getPayload())
      .then((result) => {
        if (seq !== seqRef.current) return; // stale
        setSaveStatus('saved');
        onSaved(result);
        setTimeout(() => {
          setSaveStatus((s) => (s === 'saved' ? 'idle' : s));
        }, 2000);
      })
      .catch((err: Error) => {
        if (seq !== seqRef.current) return;
        setSaveStatus('error');
        onError(err);
      });
  }, [id, getPayload, onSaved, onError, saveFn]);

  const scheduleSave = useCallback(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => {
      triggerSave();
    }, 2000);
  }, [triggerSave]);

  return { saveStatus, triggerSave, scheduleSave };
}
