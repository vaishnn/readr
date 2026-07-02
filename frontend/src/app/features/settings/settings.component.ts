import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NavbarComponent } from '../../shared/components/navbar.component';
import { ToastComponent } from '../../shared/components/toast.component';
import { SpinnerComponent } from '../../shared/components/spinner.component';
import { ToastService } from '../../shared/components/toast.service';
import { ApiService } from '../../core/services/api.service';
import { ThemeService } from '../../core/services/theme.service';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [FormsModule, NavbarComponent, ToastComponent, SpinnerComponent],
  templateUrl: './settings.component.html',
})
export class SettingsComponent {
  // Password form
  currentPassword = '';
  newPassword = '';
  confirmPassword = '';
  savingPassword = signal(false);

  constructor(
    public theme: ThemeService,
    private api: ApiService,
    private toast: ToastService,
  ) {}

  changePassword(): void {
    if (this.savingPassword()) return;

    if (this.newPassword.length < 8) {
      this.toast.error('New password must be at least 8 characters');
      return;
    }
    if (this.newPassword !== this.confirmPassword) {
      this.toast.error('Passwords do not match');
      return;
    }

    this.savingPassword.set(true);
    this.api.patch('/users/me/password', {
      currentPassword: this.currentPassword,
      newPassword: this.newPassword,
    }).subscribe({
      next: () => {
        this.toast.success('Password changed');
        this.currentPassword = '';
        this.newPassword = '';
        this.confirmPassword = '';
        this.savingPassword.set(false);
      },
      error: err => {
        this.toast.error(err?.error?.error ?? 'Failed to change password');
        this.savingPassword.set(false);
      },
    });
  }
}
