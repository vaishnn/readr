import { Component } from '@angular/core';
import { ToastService, ToastType } from './toast.service';

@Component({
  selector: 'app-toast',
  standalone: true,
  templateUrl: './toast.component.html',
})
export class ToastComponent {
  constructor(readonly toastService: ToastService) {}

  typeClass(type: ToastType): string {
    return {
      success: 'bg-emerald-800 text-emerald-100 border border-emerald-700',
      error:   'bg-red-900 text-red-100 border border-red-700',
      info:    'bg-slate-700 text-slate-100 border border-slate-600',
    }[type];
  }
}
