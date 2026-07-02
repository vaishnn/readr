export interface Collection {
  id: string;
  userId: string;
  name: string;
  description: string;
  coverUrl?: string;
  bookIds: string[];
  createdAt: string;
  updatedAt: string;
}
