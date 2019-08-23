PATH="$PATH:/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin"
TARGET_DIR=/usr/local/bin/clickhouse-gen
PERM="chmod +x /usr/local/bin/clickhouse-gen"

if [ "$(uname)" == "Darwin" ]; then
    OS="osx"
else
    OS="linux"
fi

URL="https://raw.githubusercontent.com/OrlovEvgeny/clickhouse-gen/master/build/$OS/clickhouse-gen"

if [ -n "`which curl`" ]; then
    download_cmd="curl -L $URL --output $TARGET_DIR"
else
    die "Failed to download clickhouse-gen: curl not found, plz install curl"
fi

/bin/echo "Fetching clickhouse-gen from $URL: "
$download_cmd || die "Error when downloading clickhouse-gen from $URL"
/bin/echo "Install clickhouse-gen: done"

$PERM || die "Error permission execut clickhouse-gen from $TARGET_DIR"
/bin/echo "Set permission execute clickhouse-gen: done"
/bin/echo "Finished"