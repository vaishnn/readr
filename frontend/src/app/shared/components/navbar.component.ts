import { Component, computed } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { FeatureFlagDirective } from '../../core/feature-flags';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, FeatureFlagDirective],
  templateUrl: './navbar.component.html',
})
export class NavbarComponent {
  menuOpen = false;

  private user = this.auth.currentUser;

  username = computed(() => this.user()?.username ?? '');
  email    = computed(() => this.user()?.email ?? '');
  initial  = computed(() => (this.user()?.username?.[0] ?? '?').toUpperCase());

  constructor(private auth: AuthService) {}

  logout(): void {
    this.menuOpen = false;
    this.auth.logout();
  }
}
