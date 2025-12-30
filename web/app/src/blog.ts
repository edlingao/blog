import Alpine from "alpinejs";
import htmx from "htmx.org";

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
  loading: false,
  activeReplyId: null as string | null,

  init() {
    this.$el.addEventListener("htmx:beforeRequest", (e: Event) => {
      const target = e.target as HTMLElement;
      if (!target.closest(".inline-reply-form")) {
        this.loading = true;
      }
    });
    this.$el.addEventListener("htmx:afterRequest", (e: Event) => {
      const target = e.target as HTMLElement;
      if (!target.closest(".inline-reply-form")) {
        this.resetForm();
      } else {
        this.closeReplyForm();
      }
    });
    this.$el.addEventListener("htmx:responseError", () => {
      this.loading = false;
    });
  },

  reply(commentId: string) {
    this.closeReplyForm();
    const comment = document.getElementById(`comment-${commentId}`);
    if (!comment) return;

    this.activeReplyId = commentId;
    const formHtml = this.renderReplyForm(commentId);
    comment.insertAdjacentHTML("beforeend", formHtml);

    const textarea = comment.querySelector<HTMLTextAreaElement>(
      ".reply-form-textarea",
    );
    textarea?.focus();

    const form = comment.querySelector(".inline-reply-form");
    if (form) htmx.process(form);
  },

  closeReplyForm() {
    const form = document.querySelector(".inline-reply-form");
    form?.remove();
    this.activeReplyId = null;
  },

  renderReplyForm(parentId: string): string {
    const postTitle = window.location.pathname.split("/")[2];
    return `
      <div class="inline-reply-form mt-4 aero-glass bg-white/15 backdrop-blur rounded-xl p-4">
        <form
          hx-post="/post/${postTitle}/comments"
          hx-swap="none"
          class="space-y-3"
        >
          <input type="hidden" name="reply_to" value="${parentId}"/>
          <input
            type="text"
            name="author"
            required
            placeholder="your name"
            class="w-full px-3 py-2 rounded-lg bg-white/80 border border-white/60 font-mono text-sm text-aero-deep placeholder:text-aero-deep/30 focus:outline-none focus:ring-2 focus:ring-aero-azure/40"
          />
          <textarea
            name="content"
            required
            rows="2"
            placeholder="your reply..."
            class="reply-form-textarea w-full px-3 py-2 rounded-lg bg-white/80 border border-white/60 font-body text-sm text-aero-deep placeholder:text-aero-deep/30 focus:outline-none focus:ring-2 focus:ring-aero-azure/40 resize-none"
          ></textarea>
          <div class="flex gap-2">
            <button type="submit" class="px-4 py-2 rounded-lg font-mono text-xs text-white bg-gradient-to-b from-blue-400 to-blue-600 hover:from-blue-500 hover:to-blue-700 cursor-pointer">
              reply
            </button>
            <button type="button" onclick="document.querySelector('.inline-reply-form')?.remove()" class="px-4 py-2 rounded-lg font-mono text-xs text-aero-deep/60 hover:text-red-500 cursor-pointer">
              cancel
            </button>
          </div>
        </form>
      </div>
    `;
  },

  resetForm() {
    this.content = "";
    this.loading = false;
  },
}));
