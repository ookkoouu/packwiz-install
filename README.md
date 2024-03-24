# packwiz-install
A Go port of [packwiz/packwiz-installer](https://github.com/packwiz/packwiz-installer).  
Install and update modpack with simple commands.  

# Usage
1. [Download binary from release](https://github.com/ookkoouu/packwiz-install/releases/latest) and put it on `.minecraft` folder.
2. Run install command.
```
packwiz-install install <URL>
```

## Options
Run `packwiz-install -h` for more detail.

```
$ packwiz-install install -h

Install and update modpack

Usage:
  packwiz-install install [flags] URL

Aliases:
  install, i

Flags:
  -d, --dir string    Directory to install modpack (default ".")
      --hash string   Hash of 'pack.toml' in the form of "<format>:<hash>" e.g. "sha256:abc012..."
  -h, --help          help for install
```

## Update on launch game
1. Bundle binary with your modpack.
2. Set `packwiz-install i <URL>` to Pre-Launch Hook. The hook feature is available in [Prism Launcher](https://prismlauncher.org/), [Modrinth App](https://modrinth.com/app) etc.
3. Start the game as usual.
