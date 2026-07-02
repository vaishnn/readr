import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, ActivatedRoute, RouterLink } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { ToastService } from '../../shared/components/toast.service';
import { SpinnerComponent } from '../../shared/components/spinner.component';

@Component({
  selector: 'app-auth',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, SpinnerComponent],
  templateUrl: './auth.component.html',
})
export class AuthComponent implements OnInit {
  isRegister = false;
  loading = false;
  form!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.isRegister = this.router.url.includes('register');
    this.buildForm();
  }

  private buildForm(): void {
    this.form = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      ...(this.isRegister && {
        username: ['', [Validators.required, Validators.minLength(3)]],
      }),
    });
  }

  submit(): void {
    if (this.form.invalid || this.loading) return;
    this.loading = true;

    const { email, password, username } = this.form.value;

    const request$ = this.isRegister
      ? this.auth.register(email, username, password)
      : this.auth.login(email, password);

    request$.subscribe({
      next: () => this.router.navigate(['/library']),
      error: err => {
        this.loading = false;
        const msg = err?.error?.error ?? (this.isRegister ? 'Registration failed' : 'Invalid email or password');
        this.toast.error(msg);
      },
    });
  }

  get emailError(): string {
    const c = this.form.get('email');
    if (c?.touched && c.hasError('required')) return 'Email is required';
    if (c?.touched && c.hasError('email')) return 'Enter a valid email';
    return '';
  }

  get passwordError(): string {
    const c = this.form.get('password');
    if (c?.touched && c.hasError('required')) return 'Password is required';
    if (c?.touched && c.hasError('minlength')) return 'Password must be at least 8 characters';
    return '';
  }

  get usernameError(): string {
    const c = this.form.get('username');
    if (c?.touched && c.hasError('required')) return 'Username is required';
    if (c?.touched && c.hasError('minlength')) return 'Username must be at least 3 characters';
    return '';
  }
}
