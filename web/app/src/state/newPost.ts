import Alpine from "alpinejs";

interface Metadata {
  tags: string[];
  description: string;
  comments: boolean;
}

const AVAILABLE_TAGS = [
  "books",
  "tutorials",
  "music",
  "projects",
  "art",
  "photography",
  "travel",
  "food",
  "technology",
  "gaming",
  "movies",
];

function parseFrontmatter(content: string): Metadata {
  const metadata: Metadata = {
    tags: [],
    description: "",
    comments: true,
  };

  const frontmatterMatch = content.match(/^---\s*\n([\s\S]*?)\n---/);
  if (!frontmatterMatch) return metadata;

  const frontmatter = frontmatterMatch[1];

  const tagsMatch = frontmatter.match(/tags:\s*(.+)/i);
  if (tagsMatch) {
    console.log("Tags found:", tagsMatch[1]);
    metadata.tags = tagsMatch[1]
      .replaceAll(/[\[\]]/g, "")
      .split(",")
      .map((t) => t.trim().toLowerCase())
      .filter((t) => AVAILABLE_TAGS.includes(t));
  }

  const descMatch = frontmatter.match(/description:\s*(.+)/i);
  if (descMatch) {
    metadata.description = descMatch[1].trim().replace(/^["']|["']$/g, "");
  }

  const commentsMatch = frontmatter.match(/comments:\s*(true|false)/i);
  if (commentsMatch) {
    metadata.comments = commentsMatch[1].toLowerCase() === "true";
  }

  return metadata;
}

export function NewPostForm() {
  Alpine.data("newPostForm", () => ({
    fileName: "",
    dragOver: false,
    submitting: false,
    originalTags: [] as string[],
    metadata: {
      tags: [] as string[],
      description: "",
      comments: true,
    },

    init() {
      this.clearFile();

      this.$el.addEventListener("htmx:beforeRequest", () => {
        this.submitting = true;
      });
      this.$el.addEventListener("htmx:afterRequest", () => {
        this.submitting = false;
      });
      this.$el.addEventListener("htmx:afterSwap", () => {
        this.clearFile();
      });
    },

    handleDrop(e: DragEvent) {
      this.dragOver = false;
      const files = e.dataTransfer?.files;
      if (files && files.length > 0 && files[0].name.endsWith(".md")) {
        (this.$refs.fileInput as HTMLInputElement).files = files;
        this.processFile(files[0]);
      }
    },

    handleFileSelect(e: Event) {
      const input = e.target as HTMLInputElement;
      const files = input.files;
      if (files && files.length > 0) {
        this.processFile(files[0]);
      }
    },

    processFile(file: File) {
      this.fileName = file.name;

      const titleInput = this.$refs.titleInput as HTMLInputElement;
      if (!titleInput.value.trim()) {
        titleInput.value = file.name.replace(/\.md$/i, "");
      }

      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        if (content) {
          const parsed = parseFrontmatter(content);
          this.originalTags = [...parsed.tags];
          this.metadata.tags.splice(
            0,
            this.metadata.tags.length,
            ...parsed.tags,
          );
          this.metadata.description = parsed.description;
          this.metadata.comments = parsed.comments;
        }
      };
      reader.readAsText(file);
    },

    toggleTag(tag: string) {
      if (this.originalTags.includes(tag)) return;
      const index = this.metadata.tags.indexOf(tag);
      if (index > -1) {
        this.metadata.tags.splice(index, 1);
      } else {
        this.metadata.tags.push(tag);
      }
    },

    isLockedTag(tag: string): boolean {
      return this.originalTags.includes(tag);
    },

    clearFile() {
      const fileInput = this.$refs.fileInput as HTMLInputElement;
      const titleInput = this.$refs.titleInput as HTMLInputElement;
      if (fileInput) fileInput.value = "";
      if (titleInput) titleInput.value = "";
      this.fileName = "";
      this.originalTags = [];
      this.metadata.tags.splice(0, this.metadata.tags.length);
      this.metadata.description = "";
      this.metadata.comments = true;
    },

    triggerFileInput() {
      (this.$refs.fileInput as HTMLInputElement).click();
    },
  }));
}
