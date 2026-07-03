import { Component, computed } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-shell',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './shell.component.html',
})
export class ShellComponent {
  private user = this.auth.currentUser;
  initial  = computed(() => (this.user()?.username?.[0] ?? '?').toUpperCase());
  username = computed(() => this.user()?.username ?? '');

  constructor(private auth: AuthService) {}

  logout(): void { this.auth.logout(); }
}
