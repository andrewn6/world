# https://nijmeh.cloud

Hi there! this is my personal website built with [Astro](https://astro.build)

## 📁 Folder Layout
```
/
├── public/
│   └── any assets/images
├── src/
│   ├── components/
│   │   └── any astro/react components for my website go here
│   ├── layouts/
│   │   └── the layout
│   └── pages/
│       └── the pages of my website (in .md or .astro)
└── package.json (self explanatory)
```
## 🧞 Commands

All commands are run from the root of the project, from a terminal:

*just a note I use pnpm, you may want to go install that [here](https://pnpm.io).*

| Command                | Action                                             |
| :--------------------- | :------------------------------------------------- |
| `pnpm install`          | Installs dependencies                              |
| `pnpm run dev`        | Starts local dev server at `localhost:3000`        |
| `pnpm run build`        | Build your production site to `./dist/`            |
| `pnpm run preview`      | Preview your build locally, before deploying       |
| `pnpm run astro ...`    | Run CLI commands like `astro add`, `astro preview` |
| `pnpm run astro --help` | Get help using the Astro CLI                       |
