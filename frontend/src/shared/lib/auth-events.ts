type Listener = () => void;

const listeners = new Set<Listener>();

export function subscribeAuthExpired(listener: Listener): () => void {
  listeners.add(listener);

  return () => {
    listeners.delete(listener);
  };
}

export function emitAuthExpired() {
  listeners.forEach((listener) => listener());
}
