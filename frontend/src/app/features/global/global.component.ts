import { Component, OnInit, computed, signal } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { debounceTime, distinctUntilChanged, Subject } from 'rxjs';
import { BookService } from '../../core/services/book.service';
import { Book } from '../../core/models';
import { SpinnerComponent } from '../../shared/components/spinner.component';
import { BookCardComponent } from '../library/book-card.component';

@Component({
  selector: 'app-global',
  standalone: true,
  imports: [FormsModule, SpinnerComponent, BookCardComponent],
  templateUrl: './global.component.html',
  host: { class: 'flex flex-col flex-1 min-h-0' },
})
export class GlobalComponent implements OnInit {
  books   = signal<Book[]>([]);
  total   = signal(0);
  loading = signal(true);

  searchQuery = '';
  activeTag   = '';
  page        = 1;
  readonly limit = 24;

  allTags = computed(() => [...new Set(this.books().flatMap(b => b.tags))]);

  private search$ = new Subject<string>();

  constructor(private bookService: BookService, private router: Router) {}

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
    this.bookService.list({ page: this.page, limit: this.limit, search: this.searchQuery, tag: this.activeTag })
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

  get totalPages(): number {
    return Math.ceil(this.total() / this.limit);
  }
}
