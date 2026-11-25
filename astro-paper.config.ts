import { defineAstroPaperConfig } from "./src/types/config";

export default defineAstroPaperConfig({
  site: {
    url: "https://amyasnikov.com/",
    title: "RedBlueStrategy",
    description:
      "RedBlueStrategy — notes on software engineering, system architecture and applied AI.",
    heroTitle: "Welcome",
    heroDescription:
      "Notes on software engineering, system architecture and applied AI — focused on building software that lasts.",
    author: "Alexander Myasnikov",
    profile: "https://amyasnikov.com",
    ogImage: "redbluestrategy.webp",
    lang: "en",
    timezone: "UTC",
    dir: "ltr",
  },
  posts: {
    perPage: 10,
    perIndex: 10,
  },
  features: {
    lightAndDarkMode: true,
    dynamicOgImage: true,
    showArchives: true,
    showBackButton: true,
    editPost: {
      enabled: false,
    },
    search: "pagefind",
    toc: true,
  },
  socials: [
    { name: "github",   url: "https://github.com/alexandermyasnikov" },
    { name: "linkedin", url: "https://linkedin.com/in/amyasnikov" },
    { name: "telegram", url: "https://t.me/amyasnikov" },
    { name: "gitlab",   url: "https://gitlab.com/amyasnikov" },
    { name: "mail",     url: "mailto:myasnikov.alexander.s@gmail.com" },
  ],
  shareLinks: [
    { name: "whatsapp", url: "https://wa.me/?text=" },
    { name: "facebook", url: "https://www.facebook.com/sharer.php?u=" },
    { name: "x",        url: "https://x.com/intent/post?url=" },
    { name: "telegram", url: "https://t.me/share/url?url=" },
    { name: "pinterest", url: "https://pinterest.com/pin/create/button/?url=" },
    { name: "mail",     url: "mailto:?subject=See%20this%20post&body=" },
  ],
});