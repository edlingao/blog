import "./style.css";
import Alpine from "alpinejs";
import htmx from "htmx.org";
import "./auth";
import "./blog";
import "./about";
import "./admin";
import "./toc";

htmx.config.responseHandling = [
  { code: "204", swap: false },
  { code: "[23]..", swap: true },
  { code: "[45]..", swap: true, error: true },
];

Alpine.start();
