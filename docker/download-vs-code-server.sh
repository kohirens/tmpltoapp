#!/bin/sh

# You can get the latest commit ID by looking at the latest tagged commit here: https://github.com/repos/microsoft/vscode/releases/latest
commit_id=8490d3dde47c57ba65ec40dd192d014fd2113496
hash=1612749334253
archive=vscode-server-linux-x64.tar.gz

# TODO: Find a way to get the latest commit ID via command line: https://api.github.com/repos/microsoft/vscode/releases/latest

# Download VS Code Server tarball to tmp directory.
curl -sSL "https://update.code.visualstudio.com/commit:${commit_id}/server-linux-x64/stable" -o /tmp/${archive}
# Make the parent directory where the server should live. NOTE: Ensure VS Code will have read/write access; namely the user running VScode or container user.
mkdir -p ~/.vscode-server/bin/${commit_id}_${hash}

# Extract the tarball to the right location.
tar --no-same-owner -xz --strip-components=1 -C ~/.vscode-server/bin/${commit_id}_${hash} -f /tmp/${archive}

mv -n ~/.vscode-server/bin/${commit_id}_${hash} /home/gitter/.vscode-server/bin/${commit_id}

# Make the 0 file (should add notes why this is needed).
#touch ~/.vscode-server/bin/${commit_id}/0
