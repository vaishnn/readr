import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-spinner',
  standalone: true,
  templateUrl: './spinner.component.html',
})
export class SpinnerComponent {
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
  @Input() fullPage = false;

  get sizeClass(): string {
    return { sm: 'w-4 h-4', md: 'w-8 h-8', lg: 'w-12 h-12' }[this.size];
  }

  get containerClass(): string {
    return this.fullPage ? 'fixed inset-0 bg-slate-900/60 z-50' : 'p-4';
  }
}
