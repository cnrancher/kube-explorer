#!/bin/bash

echo "basedir: $1, version: $2"
echo "to publish release to github, run push-github.sh"

echo "#!/bin/bash"> push-github.sh
echo "#"
echo "# GITHUB_TOKEN=xxx ./push-github.sh"
echo "#"
echo "pushd $1"
#echo "hub release create -t ${COMMIT} -m v${RELEASE_VERSION} v${RELEASE_VERSION}" >> dist/publish.sh
echo "for f in \$(ls $1/bin/) ; do" >> push-github.sh
echo "   hub release edit -m $2 -a $1/bin/\${f} $2" >> push-github.sh
echo "done" >> push-github.sh
echo "popd"

chmod +x push-github.sh
