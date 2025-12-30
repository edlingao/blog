import Alpine from "alpinejs";

interface CommentData {
  id: string;
  author: string;
  content: string;
  parent?: string;
}

function escapeHtml(text: string): string {
  const div = document.createElement("div");
  div.textContent = text;
  return div.innerHTML;
}

function renderComment(c: CommentData): string {
  return `<div class="aero-glass bg-white/15 backdrop-blur rounded-xl p-4 mb-4 comment-highlight" id="comment-${c.id}">
    <div class="flex items-center gap-2 mb-2">
      <span class="font-mono text-xs text-aero-azure">${escapeHtml(c.author)}</span>
    </div>
    <p class="font-body text-sm text-aero-deep/80 mb-3 leading-relaxed">${escapeHtml(c.content)}</p>
    <div class="flex items-center gap-4">
      <button
        @click="reply('${c.id}')"
        class="font-mono text-xs text-aero-deep/50 hover:text-aero-azure transition-colors cursor-pointer"
      >reply</button>
    </div>
  </div>`;
}

function createThread(parent: HTMLElement): HTMLElement {
  const thread = document.createElement("div");
  thread.className = "comment-thread mt-4";
  parent.appendChild(thread);
  return thread;
}

export function CommentEvents() {
  Alpine.data("commentEvents", () => ({
    newComments: [] as CommentData[],
    eventSource: null as EventSource | null,
    init() {
      const title = window.location.pathname.split("/")[2];
      if (!title) return;

      if (this.eventSource) {
        this.eventSource.close();
      }

      if (window.EventSource === undefined) {
        console.error("SSE not supported in this browser");
        return;
      }

      this.eventSource = new EventSource(
        `/post/${title}/stream?title=${encodeURIComponent(title)}`,
      );

      this.eventSource.addEventListener(
        "new_comment",
        (event: MessageEvent) => {
          const data = JSON.parse(event.data) as CommentData;
          this.newComments.push(data);

          let container: HTMLElement | null = null;

          if (data.parent) {
            const parent = document.getElementById(`comment-${data.parent}`);
            if (parent) {
              const thread =
                parent.querySelector<HTMLElement>(".comment-thread") ||
                createThread(parent);
              thread.insertAdjacentHTML("beforeend", renderComment(data));
              container = thread;
            }
          } else {
            const list = document.getElementById("comments-list");
            if (list) {
              const empty = list.querySelector('[class*="text-center"]');
              if (empty) empty.remove();
              list.insertAdjacentHTML("beforeend", renderComment(data));
              container = list;
            }
          }

          if (container) {
            const newComment = container.querySelector(`#comment-${data.id}`);
            if (newComment) Alpine.initTree(newComment as any);
          }
        },
      );

      this.eventSource.onerror = () => {
        console.error("SSE connection error");
      };
    },
  }));
}
