export interface Note {
  id: string;
  userId: string;
  bookId: string;
  // null means this is a book-level note; a number means it's page-specific.
  page: number | null;
  contentMd: string;
  createdAt: string;
  updatedAt: string;
}
