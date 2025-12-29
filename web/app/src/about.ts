import Alpine from "alpinejs";

Alpine.data("aboutPage", () => ({
  lang: "en",

  init() {
    const stored = localStorage.getItem("preferred_lang");
    if (stored === "en" || stored === "es") {
      this.lang = stored;
    }
  },

  toggleLang() {
    this.lang = this.lang === "en" ? "es" : "en";
    localStorage.setItem("preferred_lang", this.lang);
  },
}));
