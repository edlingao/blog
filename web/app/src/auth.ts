import Alpine from "alpinejs";

Alpine.data("authForm", () => ({
  username: "",
  password: "",
  showPassword: false,
  loading: false,

  init() {
    const form = this.$el.querySelector("form");
    if (form) {
      form.addEventListener("htmx:beforeRequest", () => (this.loading = true));
      form.addEventListener("htmx:afterRequest", () => (this.loading = false));
    }
  },
}));

Alpine.data("registerForm", () => ({
  username: "",
  email: "",
  password: "",
  confirm: "",
  showPassword: false,
  loading: false,

  get confirmError(): boolean {
    return this.confirm.length > 0 && this.password !== this.confirm;
  },

  handleSubmit(e: Event) {
    if (this.password !== this.confirm) {
      e.preventDefault();
      return;
    }
    this.loading = true;
  },
}));
