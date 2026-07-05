import { Component, Output, EventEmitter, signal, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BookService } from '../../core/services/book.service';
import { Book } from '../../core/models';
import { ToastService } from '../../shared/components/toast.service';
import { SpinnerComponent } from '../../shared/components/spinner.component';

@Component({
  selector: 'app-upload-modal',
  standalone: true,
  imports: [SpinnerComponent, FormsModule],
  templateUrl: './upload-modal.component.html',
})
export class UploadModalComponent implements OnDestroy {
  @Output() uploaded = new EventEmitter<Book>();
  @Output() close    = new EventEmitter<void>();

  bookFile        = signal<File | null>(null);
  coverFile       = signal<File | null>(null);
  thumbnailUrl    = signal<string | null>(null);
  generatingThumb = signal(false);
  uploading       = signal(false);
  dragOver        = signal(false);

  // Details form — pre-filled from PDF metadata, all editable
  bookTitle   = '';
  authorName  = '';
  tagsInput   = '';
  language    = '';
  year: number | null = null;
  description = '';

  constructor(private bookService: BookService, private toast: ToastService) {}

  ngOnDestroy(): void {
    const url = this.thumbnailUrl();
    if (url) URL.revokeObjectURL(url);
  }

  onBookFileDrop(event: DragEvent): void {
    event.preventDefault();
    this.dragOver.set(false);
    const file = event.dataTransfer?.files[0];
    if (file) this.setBookFile(file);
  }

  onBookFileSelect(event: Event): void {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (file) this.setBookFile(file);
  }

  private async setBookFile(file: File): Promise<void> {
    const ext = file.name.split('.').pop()?.toLowerCase();
    if (!['pdf', 'epub', 'cbz'].includes(ext ?? '')) {
      this.toast.error('Only PDF, EPUB and CBZ files are supported');
      return;
    }

    const prev = this.thumbnailUrl();
    if (prev) URL.revokeObjectURL(prev);
    this.thumbnailUrl.set(null);
    this.coverFile.set(null);

    // Pre-fill title from filename as baseline
    this.bookTitle   = file.name.replace(/\.[^/.]+$/, '');
    this.authorName  = '';
    this.tagsInput   = '';
    this.language    = '';
    this.year        = null;
    this.description = '';
    this.bookFile.set(file);

    if (ext === 'pdf') {
      this.generatingThumb.set(true);
      try {
        const { thumb, meta } = await this.processPdf(file);
        if (thumb) {
          this.coverFile.set(thumb);
          this.thumbnailUrl.set(URL.createObjectURL(thumb));
        }
        if (meta.title)       this.bookTitle   = meta.title;
        if (meta.author)      this.authorName  = meta.author;
        if (meta.language)    this.language    = meta.language;
        if (meta.description) this.description = meta.description;
      } finally {
        this.generatingThumb.set(false);
      }
    }
  }

  private async processPdf(file: File): Promise<{ thumb: File | null; meta: { title: string; author: string; language: string; description: string } }> {
    let thumb: File | null = null;
    const meta = { title: '', author: '', language: '', description: '' };

    try {
      const pdfjsLib = await import('pdfjs-dist');
      pdfjsLib.GlobalWorkerOptions.workerSrc = '/assets/pdf.worker.min.js';

      const arrayBuffer = await file.arrayBuffer();
      const pdf = await pdfjsLib.getDocument({ data: arrayBuffer }).promise;

      // Metadata — best effort, never breaks thumbnail generation
      try {
        const result = await pdf.getMetadata();
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const info = result?.info as any;
        meta.title       = info?.Title       ?? '';
        meta.author      = info?.Author      ?? '';
        meta.language    = info?.Language    ?? '';
        meta.description = info?.Subject     ?? '';
      } catch { /* no metadata */ }

      const page = await pdf.getPage(1);
      const viewport = page.getViewport({ scale: 1.5 });
      const canvas = document.createElement('canvas');
      canvas.width  = viewport.width;
      canvas.height = viewport.height;
      await page.render({ canvasContext: canvas.getContext('2d')!, viewport }).promise;
      pdf.destroy();

      thumb = await new Promise<File | null>(resolve => {
        canvas.toBlob(
          blob => resolve(blob ? new File([blob], 'cover.jpg', { type: 'image/jpeg' }) : null),
          'image/jpeg',
          0.88,
        );
      });
    } catch { /* thumbnail failed */ }

    return { thumb, meta };
  }

  upload(): void {
    const file = this.bookFile();
    if (!file || this.uploading()) return;
    if (!this.bookTitle.trim()) return;
    this.uploading.set(true);

    const tags = this.tagsInput.split(',').map(t => t.trim()).filter(Boolean);

    this.bookService.upload(file, this.coverFile() ?? undefined, {
      title:  this.bookTitle.trim(),
      author: this.authorName.trim() || undefined,
      tags:   tags.length ? tags : undefined,
      metadata: {
        language:    this.language.trim()    || undefined,
        year:        this.year               ?? undefined,
        description: this.description.trim() || undefined,
      },
    }).subscribe({
      next: book => {
        this.toast.success(`"${book.title}" uploaded`);
        this.uploaded.emit(book);
      },
      error: err => {
        this.uploading.set(false);
        this.toast.error(err?.error?.error ?? 'Upload failed');
      },
    });
  }

  formatSize(bytes: number): string {
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }
}
