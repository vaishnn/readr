import { Component, OnInit, computed, signal } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { debounceTime, distinctUntilChanged, Subject } from 'rxjs';
import { BookService } from '../../core/services/book.service';
import { Book } from '../../core/models';
import { SpinnerComponent } from '../../shared/components/spinner.component';
import { ToastComponent } from '../../shared/components/toast.component';
import { BookCardComponent } from './book-card.component';
import { UploadModalComponent } from './upload-modal.component';

@Component({
  selector: 'app-library',
  standalone: true,
  imports: [FormsModule, SpinnerComponent, ToastComponent, BookCardComponent, UploadModalComponent],
  templateUrl: './library.component.html',
  host: { class: 'flex flex-col flex-1 min-h-0' },
})
export class LibraryComponent implements OnInit {
  books      = signal<Book[]>([]);
  total      = signal(0);
  loading    = signal(true);
  showUpload = signal(false);

  searchQuery = '';
  activeTag   = '';
  page        = 1;
  readonly limit = 24;

  recentBooks = computed(() => this.books().slice(0, 8));

  tagCounts = computed(() => {
    const map = new Map<string, number>();
    for (const book of this.books()) {
      for (const tag of book.tags) {
        map.set(tag, (map.get(tag) ?? 0) + 1);
      }
    }
    return map;
  });

  displayTags = computed(() => {
    const tags = new Set(this.books().flatMap(b => b.tags));
    return [...tags].slice(0, 6);
  });

  private search$ = new Subject<string>();

  constructor(
    private bookService: BookService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    this.loadBooks();
    this.search$.pipe(debounceTime(300), distinctUntilChanged()).subscribe(q => {
      this.searchQuery = q;
      this.page = 1;
      this.loadBooks();
    });
  }

  loadBooks(): void {
    this.loading.set(true);
    this.bookService.list({ page: this.page, limit: this.limit, search: this.searchQuery, tag: this.activeTag, ownerOnly: true })
      .subscribe({
        next: res => {
          this.books.set(res.books ?? []);
          this.total.set(res.total ?? 0);
          this.loading.set(false);
        },
        error: () => this.loading.set(false),
      });
  }

  onSearch(query: string): void {
    this.search$.next(query);
  }

  filterByTag(tag: string): void {
    this.activeTag = tag === '' || this.activeTag === tag ? '' : tag;
    this.page = 1;
    this.loadBooks();
  }

  clearFilters(): void {
    this.activeTag = '';
    this.searchQuery = '';
    this.page = 1;
    this.loadBooks();
  }

  openBook(book: Book): void {
    this.router.navigate(['/reader', book.id]);
  }

  onBookDeleted(id: string): void {
    this.books.update(books => books.filter(b => b.id !== id));
    this.total.update(t => t - 1);
  }

  onBookUpdated(updated: Book): void {
    this.books.update(books => books.map(b => b.id === updated.id ? updated : b));
  }

  onUploaded(book: Book): void {
    this.books.update(books => [book, ...books]);
    this.total.update(t => t + 1);
    this.showUpload.set(false);
  }

  get totalPages(): number {
    return Math.ceil(this.total() / this.limit);
  }

  get allTags(): string[] {
    const tags = new Set(this.books().flatMap(b => b.tags));
    return [...tags];
  }
}
