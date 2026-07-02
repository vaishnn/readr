import { Component, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Collection, Book } from '../../core/models';
import { CollectionService } from '../../core/services/collection.service';
import { BookService } from '../../core/services/book.service';
import { ToastService } from '../../shared/components/toast.service';
import { NavbarComponent } from '../../shared/components/navbar.component';
import { SpinnerComponent } from '../../shared/components/spinner.component';
import { ToastComponent } from '../../shared/components/toast.component';

@Component({
  selector: 'app-collections',
  standalone: true,
  imports: [FormsModule, NavbarComponent, SpinnerComponent, ToastComponent],
  templateUrl: './collections.component.html',
})
export class CollectionsComponent implements OnInit {
  collections  = signal<Collection[]>([]);
  allBooks     = signal<Book[]>([]);
  loading      = signal(true);
  expandedId   = signal<string | null>(null);

  // Create form state
  showCreateForm = signal(false);
  newName        = '';
  newDescription = '';
  creating       = signal(false);

  // Edit form state
  editingId   = signal<string | null>(null);
  editName    = '';
  editDesc    = '';

  // Add book dropdown state
  addBookCollectionId = signal<string | null>(null);
  bookSearch          = '';

  constructor(
    private collectionService: CollectionService,
    private bookService: BookService,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.loadCollections();
    this.bookService.list({ limit: 200 }).subscribe({
      next: res => this.allBooks.set(res.books),
    });
  }

  loadCollections(): void {
    this.loading.set(true);
    this.collectionService.list().subscribe({
      next: cols => {
        this.collections.set(cols);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  create(): void {
    if (!this.newName.trim() || this.creating()) return;
    this.creating.set(true);
    this.collectionService.create(this.newName.trim(), this.newDescription.trim()).subscribe({
      next: col => {
        this.collections.update(cols => [col, ...cols]);
        this.newName = '';
        this.newDescription = '';
        this.showCreateForm.set(false);
        this.creating.set(false);
        this.toast.success('Collection created');
      },
      error: () => {
        this.creating.set(false);
        this.toast.error('Failed to create collection');
      },
    });
  }

  startEdit(col: Collection): void {
    this.editingId.set(col.id);
    this.editName = col.name;
    this.editDesc = col.description;
  }

  saveEdit(col: Collection): void {
    this.collectionService.update(col.id, this.editName.trim(), this.editDesc.trim()).subscribe({
      next: updated => {
        this.collections.update(cols => cols.map(c => c.id === updated.id ? updated : c));
        this.editingId.set(null);
        this.toast.success('Collection updated');
      },
      error: () => this.toast.error('Failed to update collection'),
    });
  }

  delete(id: string): void {
    if (!confirm('Delete this collection? Books will not be removed.')) return;
    this.collectionService.delete(id).subscribe({
      next: () => {
        this.collections.update(cols => cols.filter(c => c.id !== id));
        this.toast.success('Collection deleted');
      },
      error: () => this.toast.error('Failed to delete collection'),
    });
  }

  addBook(collectionId: string, bookId: string): void {
    this.collectionService.addBook(collectionId, bookId).subscribe({
      next: () => {
        this.collections.update(cols => cols.map(c =>
          c.id === collectionId && !c.bookIds.includes(bookId)
            ? { ...c, bookIds: [...c.bookIds, bookId] }
            : c
        ));
        this.addBookCollectionId.set(null);
        this.bookSearch = '';
        this.toast.success('Book added to collection');
      },
      error: () => this.toast.error('Failed to add book'),
    });
  }

  removeBook(collectionId: string, bookId: string): void {
    this.collectionService.removeBook(collectionId, bookId).subscribe({
      next: () => {
        this.collections.update(cols => cols.map(c =>
          c.id === collectionId ? { ...c, bookIds: c.bookIds.filter(id => id !== bookId) } : c
        ));
        this.toast.success('Book removed');
      },
      error: () => this.toast.error('Failed to remove book'),
    });
  }

  toggleCreateForm(): void {
    this.showCreateForm.update(v => !v);
  }

  toggleExpand(id: string): void {
    this.expandedId.update(v => v === id ? null : id);
  }

  booksInCollection(col: Collection): Book[] {
    return this.allBooks().filter(b => col.bookIds.includes(b.id));
  }

  filteredBooksToAdd(col: Collection): Book[] {
    const q = this.bookSearch.toLowerCase();
    return this.allBooks().filter(b =>
      !col.bookIds.includes(b.id) &&
      (b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q))
    );
  }
}
