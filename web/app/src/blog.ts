import Alpine from "alpinejs";

Alpine.data("postsFilter", () => ({
  activeTag: "",
  loading: false,

  init() {
    this.$el.addEventListener("htmx:beforeRequest", () => {
      this.loading = true;
    });
    this.$el.addEventListener("htmx:afterRequest", () => {
      this.loading = false;
    });
  },

  selectTag(tagId: string) {
    this.activeTag = tagId;
  },
}));

Alpine.data("commentForm", () => ({
  author: "",
  content: "",
  replyTo: null as string | null,
  loading: false,

  init() {
    this.$el.addEventListener("htmx:beforeRequest", () => {
      this.loading = true;
    });
    this.$el.addEventListener("htmx:afterSwap", () => {
      this.resetForm();
    });
    this.$el.addEventListener("htmx:responseError", () => {
      this.loading = false;
    });
  },

  reply(commentId: string) {
    this.replyTo = commentId;
    (this.$refs.commentInput as HTMLTextAreaElement)?.focus();
  },

  cancelReply() {
    this.replyTo = null;
  },

  resetForm() {
    this.content = "";
    this.replyTo = null;
    this.loading = false;
  },
}));
