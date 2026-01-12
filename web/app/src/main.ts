import "./style.css";
import Alpine from "alpinejs";
import htmx from "htmx.org";
import "./auth";
import "./blog";
import "./about";
import "./admin";
import "./toc";
import { CommentEvents } from "./state/commentEvents";
import { NewPostForm } from "./state/newPost";

import hljs from "highlight.js";
import "highlight.js/styles/github.css";

htmx.config.responseHandling = [
  { code: "204", swap: false },
  { code: "[23]..", swap: true },
  { code: "[45]..", swap: true, error: true },
];

CommentEvents();
NewPostForm();

document.addEventListener("DOMContentLoaded", () => {
  hljs.highlightAll();
  Alpine.start();
  document.body.addEventListener("htmx:afterSwap", () => {
    hljs.highlightAll();
  });
});
