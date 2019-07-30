#! /bin/bash

command="gcloud builds describe"
options="--format json --project cloudnoa-dev"

list=(
    "SUCCESS    7486540e-735e-4328-8ad8-c8c10c5fbd42"
    "CANCELLED  711f447a-ab82-493d-8e1b-0fe0105c3060"
    "TIMEOUT    b039858e-0c28-4b7a-8024-fe03c41d8403"
    "FAILURE	bd47c6c9-171c-4986-a60c-bca551fe8662"
)

mkdir -p example

for item in "${list[@]}"; do
    tokens=(${item// / })
    status=${tokens[0]}
    id=${tokens[1]}
    [ ! -s example/$status.json ] && $command $id $options > example/$status.json
done

# [ ! -s example/SUCCESS.json ] && $command 7486540e-735e-4328-8ad8-c8c10c5fbd42 $options > example/SUCCESS.json
# [ ! -s example/CANCELLED.json ] && $command 711f447a-ab82-493d-8e1b-0fe0105c3060 $options > example/CANCELLED.json
