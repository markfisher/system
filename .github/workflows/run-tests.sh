#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

source ${FATS_DIR}/.configure.sh
source ${FATS_DIR}/applications/helpers.sh
source ${FATS_DIR}/functions/helpers.sh

# setup namespace
kubectl create namespace ${NAMESPACE}
fats_create_push_credentials ${NAMESPACE}

for test in function application; do
  path="${test}s/uppercase/node"
  name=system-${RUNTIME}-uppercase-node-${test}
  image=$(fats_image_repo ${name})
  create_args="--git-repo https://github.com/${FATS_REPO}.git --git-revision ${FATS_REFSPEC} --sub-path ${path}"
  input_data=system
  expected_data=SYSTEM
  runtime=$RUNTIME

  run_${test} ${FATS_DIR}/$path $name $image "$create_args" $input_data $expected_data $runtime
done
