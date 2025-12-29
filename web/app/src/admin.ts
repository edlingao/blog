import Alpine from "alpinejs";

Alpine.data("adminPosts", () => ({
  confirmDelete: null as string | null,
  refreshing: false,

  init() {
    this.$el.addEventListener("htmx:beforeRequest", (e: Event) => {
      const detail = (e as CustomEvent).detail;
      if (detail?.path?.includes("refresh")) {
        this.refreshing = true;
      }
    });
    this.$el.addEventListener("htmx:afterRequest", (e: Event) => {
      const detail = (e as CustomEvent).detail;
      if (detail?.path?.includes("refresh")) {
        this.refreshing = false;
      }
    });
    this.$el.addEventListener("htmx:afterSwap", () => {
      this.confirmDelete = null;
    });
  },

  askDelete(postId: string) {
    this.confirmDelete = postId;
  },

  cancelDelete() {
    this.confirmDelete = null;
  },
}));

Alpine.data("adminComments", () => ({
  confirmDelete: null as string | null,

  init() {
    this.$el.addEventListener("htmx:afterSwap", () => {
      this.confirmDelete = null;
    });
  },

  askDelete(commentId: string) {
    this.confirmDelete = commentId;
  },

  cancelDelete() {
    this.confirmDelete = null;
  },
}));
