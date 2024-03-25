# packwiz-install
An standalone updater for packwiz modpacks with simple commands. No need to install Java anymore!  
I hope it replaces [packwiz-installer](https://github.com/packwiz/packwiz-installer) and [packwiz-installer-bootstrap](https://github.com/packwiz/packwiz-installer-bootstrap).  

## Usage
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
2. Set Pre-Launch Hook to player's launcher. The hook feature is available in [Prism Launcher](https://prismlauncher.org/), [Modrinth App](https://modrinth.com/app) etc.
    ![image](https://github.com/ookkoouu/packwiz-install/assets/29059223/cd22914c-d09a-46e6-8852-bf3cff78fc3e)
    ### Windows
    ```
    cmd /c packwiz-install install <URL>
    ```
    ### Linux
    ```
    sh -c 'packwiz-install install <URL>'
    ```
4. Start the game as usual.
