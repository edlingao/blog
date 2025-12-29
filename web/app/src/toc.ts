import Alpine from "alpinejs";

Alpine.data("tableOfContents", () => ({
  open: true,
  hasToc: false,

  init() {
    this.$nextTick(() => {
      this.moveTocToPanel();
    });
  },

  moveTocToPanel() {
    const tocList = document.getElementById("TableOfContents");
    const tocTitle = document.getElementById("table-of-contents");
    const panelDesktop = document.getElementById("toc-panel-content");
    const panelMobile = document.getElementById("toc-panel-content-mobile");
    const blogContent = document.querySelector(".blog-content");

    if (tocList && (panelDesktop || panelMobile)) {
      this.hasToc = true;

      // Clone for desktop
      if (panelDesktop) {
        const cloneDesktop = tocList.cloneNode(true) as HTMLElement;
        cloneDesktop.removeAttribute("id");
        cloneDesktop.classList.add("toc-floating");
        panelDesktop.appendChild(cloneDesktop);
        this.addSmoothScroll(cloneDesktop, blogContent);
      }

      // Clone for mobile
      if (panelMobile) {
        const cloneMobile = tocList.cloneNode(true) as HTMLElement;
        cloneMobile.removeAttribute("id");
        cloneMobile.classList.add("toc-floating");
        panelMobile.appendChild(cloneMobile);
        this.addSmoothScroll(cloneMobile, blogContent);
      }

      // Hide original TOC section
      if (tocTitle) tocTitle.style.display = "none";
      tocList.style.display = "none";
    }
  },

  addSmoothScroll(container: HTMLElement, scrollContainer: Element | null) {
    const links = container.querySelectorAll("a[href^='#']");
    links.forEach((link) => {
      link.addEventListener("click", (e) => {
        e.preventDefault();
        const href = link.getAttribute("href");
        if (!href) return;

        const targetId = href.slice(1);
        const target = document.getElementById(targetId);

        if (target && scrollContainer) {
          const containerRect = scrollContainer.getBoundingClientRect();
          const targetRect = target.getBoundingClientRect();
          const offset = targetRect.top - containerRect.top + scrollContainer.scrollTop - 20;

          scrollContainer.scrollTo({
            top: offset,
            behavior: "smooth",
          });
        }
      });
    });
  },

  toggle() {
    this.open = !this.open;
  },
}));
