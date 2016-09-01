#!/bin/bash

set -e

# data and auth files
SCRIPTS_FOLDER=${HOME}/scripts
AUTH_FILE=${SCRIPTS_FOLDER}/.auth
DATA_FILE=${SCRIPTS_FOLDER}/data

# repository configuration
REMOTE_GATECH=gatech

# log
LOG_FILE=${HOME}/logs/overleaf_sync.log
exec 1<&-
exec 1<>${LOG_FILE}

# local folder configuration
OVERLEAF=${HOME}/overleaf
GATECH=${HOME}/gatech
TRANSIENT=${HOME}/transient

function handle_repo() {
  # First we create a backup of the overleaf repo
  OF_FOLDER=${OVERLEAF}/${OVERLEAF_INDEX}
  if [[ ! -d ${OF_FOLDER} ]]; then
    printf "### cloning from overleaf for backup ...\n"
    cd ${OVERLEAF} && git clone https://git.overleaf.com/${OVERLEAF_INDEX}
  fi
  printf "### pulling from overleaf for backup ...\n"
  cd ${OF_FOLDER} && git pull origin master
  unset OF_FOLDER

  # Now we create a backup of the gatech repo
  GT_FOLDER=${GATECH}/${GITHUB_REPO}
  if [[ ! -d ${GT_FOLDER} ]]; then
    printf "### cloning from gatech github for backup ...\n"
    mkdir -p ${GT_FOLDER}
    cd ${GT_FOLDER} && git init && git remote add origin https://${USERNAME}:${ACCESS_TOKEN}@github.gatech.edu/penguins/${GITHUB_REPO}.git
  else
    printf "### pulling from gatech github for backup ...\n"
    cd ${GT_FOLDER} && git pull origin master
  fi
  unset GT_FOLDER

  # clone the repository
  REPO_FOLDER=${TRANSIENT}/${OVERLEAF_INDEX}
  if [[ ! -d ${REPO_FOLDER} ]]; then
    printf "### cloning from overleaf for sync ...\n"
    cd ${TRANSIENT} && git clone https://git.overleaf.com/${OVERLEAF_INDEX}
  fi
  printf "### pulling from overleaf for sync ...\n"
  cd ${REPO_FOLDER} && git pull origin master

  # add remote
  GATECH_EXISTS=$(git remote -v | grep ${REMOTE_GATECH} | wc -l)
  if [[ "${GATECH_EXISTS}" != "2" ]]; then
    printf "### adding remote to overleaf repo ...\n"
    cd ${REPO_FOLDER} && git remote add ${REMOTE_GATECH} https://${USERNAME}:${ACCESS_TOKEN}@github.gatech.edu/penguins/${GITHUB_REPO}.git
  fi

  # push the changes
  printf "### sync overleaf with gatech github ...\n"
  cd ${REPO_FOLDER} && git push ${REMOTE_GATECH} master
}

# main
source ${AUTH_FILE}

# log
printf "\n\n$(date)\n"

# read data file
while IFS=":" read -r KEY VALUE; do
  printf "# sync : $KEY => $VALUE\n"
  OVERLEAF_INDEX=$KEY
  GITHUB_REPO=$VALUE
  handle_repo
  printf "# sync complete\n\n"
done < "${DATA_FILE}"
