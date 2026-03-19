# aidfm

**A**pp**I**mage **D**esktop **F**ile **M**anager is a CLI for managing AppImages on Linux. It automates the creation of desktop entries, manages binary permissions, and maintains a local registry to track the health of installed applications.

## Getting started

To use aidfm, build the application from source:

```bash
git clone https://github.com/yourname/aidfm.git
cd aidfm
go build -o aidfm .
mv aidfm ~/.local/bin/aidfm
```

The simplest usage would be adding an AppImage or a directory with an AppImage:

```bash
aidfm some_app.AppImage
aidfm add some_dir/some_app.AppImage
```

If presented with a directory, it will autodetect if an image file exists and will set it as the desktop icon.

## Usage

The aidfm CLI has the following structure:

```
aidfm COMMAND [ARGS] [FLAGS]
```

### AppImage management

```
add <path>                 Register an AppImage and create a .desktop file
  --global                 Install to /usr/share/applications instead of
                           ~/.local/share/applications--global

fix [name]                 Re-apply setup for a managed entry, or all broken
                           entries if no name given

remove <name>              Remove a managed entry and its desktop file
  --purge, -p              Also move the app directory to trash
  --yes, -y                Skip confirmation prompt
```

### Desktop file management

```
import <name|path>         Take ownership of an existing desktop file
edit <name>                Open the desktop file in $EDITOR
show <name>                Pretty-print all fields of a managed entry

env set <name> KEY VALUE   Set an environment variable on the Exec= line
env unset <name> KEY       Remove an environment variable from the Exec= line
env list <name>            List all environment variables on the Exec= line
```

### Registry

```
list                       List all managed entries with status
sync                       Reconcile registry against disk state
status                     Print a summary of managed, broken, and orphaned entries
```

### General
```
--help, -h                 Display help for any command
--version, -v              Display aidfm version
```

