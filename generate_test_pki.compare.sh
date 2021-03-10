#!/bin/bash

# Exits on any error
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

####################################################################
## Generate PKI for Tests
####################################################################
##  This script allows for the easy creation of a PKI via OpenSSL.
##  You can use this to instantiate your own basic PKI quickly,
##  however this is used primarily to test the Golang generated
##  PKI against a known working PKI generated via OpenSSL

function splitLinesFront {  
  echo "$1" | while IFS= read -r line ; do
    if [[ $line != *"Certificate purposes:"* ]]; then
      #echo $line
      readarray -d : -t strarr <<< "$line"
      echo "${strarr[0]}"
    fi
  done
}
function splitLinesBack {  
  echo "$1" | while IFS= read -r line ; do
    if [[ $line != *"Certificate purposes:"* ]]; then
      #echo $line
      readarray -d : -t strarr <<< "$line"
      echo "${strarr[1]}"
    fi
  done
}

function splitPurposes {
  DEF_LINES=()
  DEF_ARRVAR=()
  OSSL_LINES=()
  OSSL_ARRVAR=()
  LOCK_LINES=()
  LOCK_ARRVAR=()

  OSSL_PURPOSE="$1"
  LOCK_PURPOSE="$2"
  
  DEF_LINES=$(splitLinesFront "$OSSL_PURPOSE")
  OSSL_LINES=$(splitLinesBack "$OSSL_PURPOSE")
  LOCK_LINES=$(splitLinesBack "$LOCK_PURPOSE")

  echo "$DEF_LINES" | {
    while IFS= read -r ROW; do
      if [[ $ROW != "" ]]; then
        DEF_ARRVAR+=("$ROW")
      fi
    done
    echo "$OSSL_LINES" | {
      while IFS= read -r OROW; do
        if [[ $OROW != "" ]]; then
          OSSL_ARRVAR+=("$OROW")
        fi
      done
      
      echo "$LOCK_LINES" | {
        while IFS= read -r LROW; do
          if [[ $LROW != "" ]]; then
            LOCK_ARRVAR+=("$LROW")
          fi
        done
        
        COUNTER_INDX=0
        echo -e "------------------------------------------OSSL---LOCK-"
        for i in "${DEF_ARRVAR[@]}"
        do
          line='                                                  '
          setOSSL_D="${OSSL_ARRVAR[$COUNTER_INDX]}"
          COLOR=$RED
          if [[ "$setOSSL_D" == "${LOCK_ARRVAR[$COUNTER_INDX]}" ]]; then
            COLOR=$GREEN
          fi
          printf "${COLOR}%-40s %-6s %-6s${NC}" "$i" "$setOSSL_D" "${LOCK_ARRVAR[$COUNTER_INDX]}"
          echo -e "\n------------------------------------------------------"
          #printf "%s %s %s\n" "${line:${#i}}" "$setOSSL_D" "${LOCK_ARRVAR[$COUNTER_INDX]}"
          #echo -e "$i:${setOSSL_D} :${LOCK_ARRVAR[$COUNTER_INDX]}"
          let COUNTER_INDX=COUNTER_INDX+1
        done
      }
    }
  }

}

function checkForProgramAndExit() {
    command -v $1
    if [[ $? -eq 0 ]]; then
        printf '%-72s %-7s\n' $1 "PASSED!";
    else
        printf '%-72s %-7s\n' $1 "FAILED!";
        exit 1
    fi
}

checkForProgramAndExit openssl

####################################################################
## Set up variables

CUR_DIR=$(pwd)

OPENSSL_PKI_ROOT_DIR="${CUR_DIR}/.test_pki_root"
LOCKSMITH_PKI_ROOT_DIR="${CUR_DIR}/.generated/roots/example-labs-root-certificate-authority"

echo -e "\n===== ISSUER COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -issuer)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -issuer)"

echo -e "\n===== SUBJECT COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -subject)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -subject)"

echo -e "\n===== STARTDATE COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -startdate)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -startdate)"

echo -e "\n===== ENDDATE COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -enddate)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -enddate)"

echo -e "\n===== SERIAL COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -serial)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -serial)"

echo -e "\n===== EMAIL COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -email)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -email)"

echo -e "\n===== PURPOSE COMPARISON\n"
OSSP_PURPOSE_CMD=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -purpose)
LOCK_PURPOSE_CMD=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -purpose)

splitPurposes "$OSSP_PURPOSE_CMD" "$LOCK_PURPOSE_CMD"