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
      #readarray -d : -t strarr <<< "$line"
      #echo "${strarr[0]}"

      part=`echo $line | awk -F":" '{print $1}'`
      #readarray -d : -t strarr <<< "$line"
      echo "$part"
    fi
  done
}
function splitLinesBack {  
  echo "$1" | while IFS= read -r line ; do
    if [[ $line != *"Certificate purposes:"* ]]; then
      #echo $line
      part=`echo $line | awk -F":" '{print $2}'`
      #readarray -d : -t strarr <<< "$line"
      echo "$part"
    fi
  done
}

function splitTextOutput {
  O_ARRVAR=()
  L_ARRVAR=()
  echo "$1" | {
    while IFS= read -r ROW; do
      if [[ $ROW != "" ]]; then
        O_ARRVAR+=("$ROW")
      fi
    done
    echo "$2" | {
      while IFS= read -r LROW; do
        if [[ $LROW != "" ]]; then
          L_ARRVAR+=("$LROW")
        fi
      done
      
        
      COUNTER_INDX=0
        echo -e "----------------------------------------------OSSL--------------------------------------------------------------------------LOCK------------------------------------"
      for i in "${O_ARRVAR[@]}"
      do
        printf "%-100s %-100s\n" "$i" "${L_ARRVAR[$COUNTER_INDX]}"
        let COUNTER_INDX=COUNTER_INDX+1
      done
      
    }
  }
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
function simpleCompare {
  COLOR=$RED
  if [[ "$1" == "$2" ]]; then
    COLOR=$GREEN
  fi
  printf "${COLOR}OSSL: %s${NC}\n" "$1"
  printf "${COLOR}Lock: %s${NC}\n" "$2"
}
function simpleCompareNoColor {
  printf "OSSL: %s\n" "$1"
  printf "Lock: %s\n" "$2"
}
checkForProgramAndExit openssl

####################################################################
## Set up variables

CUR_DIR=$(pwd)

OPENSSL_PKI_ROOT_DIR="${CUR_DIR}/.test_pki_root"
LOCKSMITH_PKI_ROOT_DIR="${CUR_DIR}/.generated/roots/example-labs-root-certificate-authority"
OPENSSL_CA_CERT="ca.cert.pem"
LOCKSMITH_CA_CERT="certs/ca.pem"

echo -e "\n===== ISSUER COMPARISON\n"
OSSL_ISSUER=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -issuer)
LOCK_ISSUER=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -issuer)
simpleCompare "$OSSL_ISSUER" "$LOCK_ISSUER"

echo -e "\n===== SUBJECT COMPARISON\n"
OSSL_SUBJECT=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -subject)
LOCK_SUBJECT=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -subject)
simpleCompare "$OSSL_SUBJECT" "$LOCK_SUBJECT"

echo -e "\n===== STARTDATE COMPARISON\n"
OSSL_STARTDATE=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -startdate)
LOCK_STARTDATE=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -startdate)
simpleCompareNoColor "$OSSL_STARTDATE" "$LOCK_STARTDATE"

echo -e "\n===== ENDDATE COMPARISON\n"
OSSL_ENDDATE=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -enddate)
LOCK_ENDDATE=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -enddate)
simpleCompareNoColor "$OSSL_ENDDATE" "$LOCK_ENDDATE"

echo -e "\n===== SERIAL COMPARISON\n"
OSSL_SERIAL=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -serial)
LOCK_SERIAL=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -serial)
simpleCompare "$OSSL_SERIAL" "$LOCK_SERIAL"

echo -e "\n===== EMAIL COMPARISON\n"
OSSL_EMAIL=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -email)
LOCK_EMAIL=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -email)
simpleCompare "$OSSL_EMAIL" "$LOCK_EMAIL"

echo -e "\n===== PURPOSE COMPARISON\n"
OSSP_PURPOSE_CMD=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -purpose)
LOCK_PURPOSE_CMD=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -purpose)
splitPurposes "$OSSP_PURPOSE_CMD" "$LOCK_PURPOSE_CMD"

OSSP_TEXT_CMD=$(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/${OPENSSL_CA_CERT} -noout -text)
LOCK_TEXT_CMD=$(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/${LOCKSMITH_CA_CERT} -noout -text)
splitTextOutput "${OSSP_TEXT_CMD}" "${LOCK_TEXT_CMD}"