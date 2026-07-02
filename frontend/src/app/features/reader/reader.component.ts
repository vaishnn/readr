import { Component, OnInit, OnDestroy, signal, HostListener, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs';
import { Book, Highlight, ReadingProgress } from '../../core/models';
import { BookService } from '../../core/services/book.service';
import { HighlightService } from '../../core/services/highlight.service';
import { ToastService } from '../../shared/components/toast.service';
import { ToastComponent } from '../../shared/components/toast.component';
import { SpinnerComponent } from '../../shared/components/spinner.component';
import { PdfViewerComponent, TextSelection } from './pdf-viewer/pdf-viewer.component';
import { EpubViewerComponent } from './epub-viewer/epub-viewer.component';
import { HighlightToolbarComponent } from './highlight-toolbar/highlight-toolbar.component';
import { NotesPanelComponent } from './notes-panel/notes-panel.component';
import { HighlightColor } from '../../core/models';

@Component({
  selector: 'app-reader',
  standalone: true,
  imports: [
    RouterLink,
    SpinnerComponent,
    ToastComponent,
    PdfViewerComponent,
    EpubViewerComponent,
    HighlightToolbarComponent,
    NotesPanelComponent,
  ],
  templateUrl: './reader.component.html',
})
export class ReaderComponent implements OnInit, OnDestroy {
  @ViewChild('pdfViewer') pdfViewer?: PdfViewerComponent;

  book        = signal<Book | null>(null);
  progress    = signal<ReadingProgress | null>(null);
  highlights  = signal<Highlight[]>([]);
  loading     = signal(true);

  currentPage  = signal(1);
  currentCFI   = signal('');
  currentZoom  = signal(0); // 0 = use auto-fit default
  zoomPercent  = signal(100);
  notesOpen    = signal(false);

  selection   = signal<TextSelection | null>(null);
  toolbarPos  = signal<{ x: number; y: number } | null>(null);

  private sessionStart = Date.now();
  private saveTimer?: ReturnType<typeof setInterval>;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private bookService: BookService,
    private highlightService: HighlightService,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    const bookId = this.route.snapshot.paramMap.get('bookId')!;

    forkJoin({
      book:     this.bookService.get(bookId),
      progress: this.bookService.getProgress(bookId),
    }).subscribe({
      next: ({ book, progress }) => {
        this.book.set(book);
        if (progress) {
          this.progress.set(progress);
          this.currentPage.set(progress.page || 1);
          this.currentCFI.set(progress.cfi || '');
          this.currentZoom.set(progress.zoom || 0);
        }
        this.loadHighlights(bookId);
        this.loading.set(false);
      },
      error: () => {
        this.toast.error('Failed to load book');
        this.router.navigate(['/library']);
      },
    });

    // Auto-save progress every 30 seconds.
    this.saveTimer = setInterval(() => this.saveProgress(), 30_000);
  }

  ngOnDestroy(): void {
    clearInterval(this.saveTimer);
    this.saveProgress();
  }

  @HostListener('window:beforeunload')
  onBeforeUnload(): void {
    this.saveProgress();
  }

  // Toggle notes panel with N key.
  @HostListener('document:keydown.n')
  toggleNotes(): void {
    this.notesOpen.update(v => !v);
  }

  @HostListener('document:keydown.escape')
  onEscape(): void {
    if (this.toolbarPos()) {
      this.clearSelection();
    } else {
      this.notesOpen.set(false);
    }
  }

  onPageChange(page: number): void {
    this.currentPage.set(page);
    this.loadHighlights(this.book()!.id);
  }

  onCFIChange(cfi: string): void {
    this.currentCFI.set(cfi);
  }

  onZoomPercentChange(percent: number): void {
    this.zoomPercent.set(percent);
  }

  zoomIn(): void  { this.pdfViewer?.zoomIn(); }
  zoomOut(): void { this.pdfViewer?.zoomOut(); }

  onTextSelected(event: TextSelection): void {
    this.selection.set(event);
    this.toolbarPos.set({ x: event.x, y: event.y });
  }

  clearSelection(): void {
    this.selection.set(null);
    this.toolbarPos.set(null);
    window.getSelection()?.removeAllRanges();
  }

  onHighlightUpdated(updated: Highlight): void {
    this.highlights.update(list => list.map(h => h.id === updated.id ? updated : h));
  }

  onHighlightDeleted(id: string): void {
    this.highlights.update(list => list.filter(h => h.id !== id));
  }

  onHighlightCreate(color: HighlightColor): void {
    const sel = this.selection();
    const book = this.book();
    if (!sel || !book) return;

    this.highlightService.create(book.id, {
      page:     this.currentPage(),
      cfiRange: sel.cfiRange ?? '',
      text:     sel.text,
      color,
    }).subscribe({
      next: h => {
        this.highlights.update(list => [...list, h]);
        this.toast.success('Highlight saved');
      },
      error: () => this.toast.error('Failed to save highlight'),
    });

    this.clearSelection();
  }

  get streamUrl(): string {
    return this.bookService.streamUrl(this.book()!.id);
  }

  private loadHighlights(bookId: string): void {
    const page = this.book()?.format === 'epub' ? undefined : this.currentPage();
    this.highlightService.list(bookId, page).subscribe({
      next: h => this.highlights.set(h),
    });
  }

  private saveProgress(): void {
    const book = this.book();
    if (!book) return;
    const sessionSeconds = Math.floor((Date.now() - this.sessionStart) / 1000);
    this.bookService.saveProgress(book.id, this.currentPage(), this.currentCFI(), 0, 0, sessionSeconds).subscribe();
    this.sessionStart = Date.now();
  }
}
