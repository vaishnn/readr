export type HighlightColor = 'yellow' | 'green' | 'blue' | 'pink';

export interface Highlight {
  id: string;
  userId: string;
  bookId: string;
  page: number;
  cfiRange: string;
  text: string;
  color: HighlightColor;
  note: string;
  createdAt: string;
  updatedAt: string;
}

// Snapshot is the full highlight document captured before each change.
export interface HighlightHistory {
  id: string;
  highlightId: string;
  action: 'create' | 'update' | 'delete';
  snapshot: Record<string, unknown>;
  timestamp: string;
}
