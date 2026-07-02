import { Component, Input, Output, EventEmitter, signal } from '@angular/core';
import { Book } from '../../core/models';
import { BookService } from '../../core/services/book.service';
import { ToastService } from '../../shared/components/toast.service';

@Component({
  selector: 'app-book-card',
  standalone: true,
  templateUrl: './book-card.component.html',
})
export class BookCardComponent {
  @Input({ required: true }) book!: Book;
  @Output() open    = new EventEmitter<Book>();
  @Output() deleted = new EventEmitter<string>();

  menuOpen = signal(false);
  deleting = signal(false);

  constructor(private bookService: BookService, private toast: ToastService) {}

  get formatBadgeClass(): string {
    return {
      pdf:  'bg-red-900/60 text-red-300',
      epub: 'bg-blue-900/60 text-blue-300',
      cbz:  'bg-purple-900/60 text-purple-300',
    }[this.book.format] ?? 'bg-slate-700 text-slate-300';
  }

  delete(event: MouseEvent): void {
    event.stopPropagation();
    if (!confirm(`Delete "${this.book.title}"?`)) return;
    this.deleting.set(true);
    this.bookService.delete(this.book.id).subscribe({
      next: () => {
        this.toast.success('Book deleted');
        this.deleted.emit(this.book.id);
      },
      error: err => {
        this.deleting.set(false);
        this.toast.error(err?.error?.error ?? 'Failed to delete book');
      },
    });
  }

  toggleMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.menuOpen.update(v => !v);
  }
}
