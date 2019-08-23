PATH="$PATH:/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin"
TARGET_DIR=/usr/local/bin/clichouse-gen
PERM="chmod +x /usr/local/bin/clichouse-gen"

if [ "$(uname)" == "Darwin" ]; then
    OS="osx"
else
    OS="linux"
fi

URL="https://raw.githubusercontent.com/OrlovEvgeny/clickhouse-gen/master/build/$OS/clickhouse-gen"

if [ -n "`which curl`" ]; then
    download_cmd="curl -L $URL --output $TARGET_DIR"
else
    die "Failed to download clichouse-gen: curl not found, plz install curl"
fi

/bin/echo "Fetching clichouse-gen from $URL: "
$download_cmd || die "Error when downloading clichouse-gen from $URL"
/bin/echo "Install clichouse-gen: done"

$PERM || die "Error permission execut clichouse-gen from $TARGET_DIR"
/bin/echo "Set permission execute clichouse-gen: done"
/bin/echo "Finished"