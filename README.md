# gowl
This is a failed attempt at writing a basic copy/paste utility for [Wayland](https://wayland.freedesktop.org/).

I wanted a pure code way of doing copy/paste in a Go program without needing a CLI tool like [`wl-clipboard`](https://github.com/bugaevc/wl-clipboard). It turns out, Wayland is a big thing and I didn't have the patience to learn all the bits required to do (what I thought was) a simple thing.

`archive/` and `archive2` (because what's version control?) contain attempts with the [`go-wayland`](https://github.com/rajveermalviya/go-wayland) and [`go-wlroots`](https://github.com/swaywm/go-wlroots) libraries, respectively. The stuff in the current directory tries the [`neurlang/wayland`](https://github.com/neurlang/wayland) library.

I did learn about the Go `generate` command. I'll take that small victory and leave this project behind.
