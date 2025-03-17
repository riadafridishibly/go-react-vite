## Go + Vite + React


This repo is some kind of template to initilalize react (vite) project with golnag using (gin) framework. This is quite simple, but it gets complicated when serving (react) assets file. For example we have `/api` routes and we have UI routes which is used by tanstack router. Everything works great until user reloads a page. So the trick is to serve `index.html` when the route is not an API route. 