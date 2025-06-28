#!/bin/bash

while IFS='=' read -r key value; do
    [[ $key ]] || continue  # Skip empty lines
    key=$(echo "$key" | tr -d ' ')  # Remove spaces from key
    value=$(echo "$value" | sed 's/^["'\'']//;s/["'\'']$//')  # Remove surrounding quotes
    export "$key=$value"
done < <(grep -v '^#' .env)


dbName="${DB_XI:-XI}"
confName="xi.com.4000"

dbPass="${DB_XI_PASS:-${DB_PASS:-"noPass"}}"
dbSrcDir="asset/space/setup/db"
dbSrc="${dbSrcDir}/${dbName}"
dbHost=${DB_XI_HOST:-localhost}

confSrcDir="asset/space/setup/conf"
confSrc="${confSrcDir}/${confName}"

confDestDir="/etc/nginx/sites-available";
[[ -n "$PREFIX" ]] && confDestDir="${PREFIX}${confDestDir}"
confDest="${confDestDir}/${confName}"
confEnDir="$(realpath "${confDestDir}/../sites-enabled/")"
confEn="$(realpath "${confEnDir}/${confName}")"

composerCacheDir=/var/www/.cache/; [[ -n "$PREFIX" ]] && composerCacheDir="${PREFIX}${composerCacheDir}"

safeUser=www-data;
if ! id "www-data" &>/dev/null; then 
    safeUser="${SUDO_USER:-${LOGNAME:-nobody}}";
    
    [[ "$safeUser" == "root" || -z "$safeUser" ]] && safeUser="nobody"
    
    # if id "everybody" &>/dev/null; then 
    #     safeUser="everybody";
    # fi
fi

# Dynamic val
PHP_VER=$(php -r 'echo PHP_MAJOR_VERSION.".".PHP_MINOR_VERSION;')


scriptPath="$(readlink -f ${BASH_SOURCE[0]:-$0})"
scriptDir="$(dirname "${scriptPath}")";
scriptName=${0##*/}

# Initialize an empty array for color history
cc=()
out() {
    local c="$1"; local msg="$2"; local mode="$3"; local cx="$4"
    
    msg="$(sed "s/--\\\\c/${c//\\/\\\\}/g" <<< "$msg")"
    
    # Push new color to front of stack
    cc=("$c" "${cc[@]}")
    # echo -e "\nc3=${cc[3]:-} c2=${cc[2]:-} c1=${cc[1]:-} c0=${cc[0]}\n"
    local x="${cc[$cx]:-\e[0m}"

    case "$mode" in
        n) printf "${c}${msg}\e[0m" ;;
        p) printf "${c}${msg}${x}" ;;
        pc) printf "${x}${msg}${x}" ;;
        pn) printf "${c}${msg}${x}\n" ;;
        *) printf "${c}${msg}\e[0m\n" ;;
    esac
}
pout() { out "${2}" "${1}" "p" "$3"; }
msg() { out "\e[0;36m" "${1}" "${2}" "$3"; }
smsg() { out "\e[0;36m" "[${scope:-$scriptName}]\e[0;2m ${1}" "${2}" "$3"; }

inf() { out "\e[0;2m" "${1}" "${2}" "$3"; }
sinf() { out "\e[0;2m" "\e[0;36m[${scope:-$scriptName}] --\c${1}" "${2}" "$3"; }
sinf1() { out "\e[0;2m" "[\e[0;36m${scope:-$scriptName}--\c] ${1}" "${2}" "$3"; }

alrt() { out "\e[0;35m" "${1}" "${2}" "$3"; }
salrt() { out "\e[0;35m" "\e[0;36m[${scope:-$scriptName}] --\c${1}" "${2}" "$3"; }
salrt1() { out "\e[0;35m" "[\e[0;36m${scope:-$scriptName}--\c ${1}" "${2}" "$3"; }

wrn() { out "\e[0;33m" "${1} !!" "${2}" "$3"; return 1; }
swrn() { out "\e[0;33m" "\e[0;36m[${scope:-$scriptName}] --\c${1}" "${2}" "$3"; }
swrn1() { out "\e[0;33m" "[\e[0;36m${scope:-$scriptName}--\c ${1}" "${2}" "$3"; }

err() { out "\e[0;31m" "${1} !!" "${2}" "$3"; exit 1; }
serr() { out "\e[0;31m" "\e[0;36m[${scope:-$scriptName}] --\c${1}" "${2}" "$3"; }
serr1() { out "\e[0;31m" "[\e[0;36m${scope:-$scriptName}--\c] ${1}" "${2}" "$3"; }



[ -f /x/bin/sh/lib-base-sh ] && source /x/bin/sh/lib-base-sh

function fix(){
    scope=Fix;

    if [[ -e ./asset/space/setup/setup.sh ]]; then
        if [ "$(id -u)" -eq 0 ]; then
            u2=", '${safeUser}'"
        fi
        smsg "(git) $(inf "Configuring safe dir for ['$(whoami)'$u2] to '$PWD'")"
        git config --global --add safe.directory $PWD

        if [ "$(id -u)" -eq 0 ]; then
            sudo -u ${safeUser} git config --global --add safe.directory "$PWD"
            smsg "(owner) $(inf "Configuring ['${safeUser}:${safeUser}:775'] '$PWD'")"
            chown ${safeUser}:${safeUser} "$PWD" -R
            chmod 775 ./ -R
        fi

    fi
}

function _pkg(){
    scope=Pkg;

    function setup(){
        msg "[Pkg] $(inf "Installing required packages")"
        apt install -y mariadb
        apt install -y mariadb-server mariadb-client
        apt install -y git nginx php php-fpm composer curl unzip openssl
        apt install -y php-xml php-sqlite3 php-mysql
        apt install -y keepalived avahi-daemon avahi-utils
    }

    case "$1" in
        install|setup) setup
            shift ;;
        -b|b|bkp|backup) backup
            shift ;;
        u|update) update
            shift ;;
        :) echo "Usage: pkg {{install|setup}|update}" ;;
    esac
}

function _nginx(){
    scope=Nginx;
    if [ "${nginxSig:-1}" -eq 0 ]; then
        sinf "$(alrt "Skipping nginx related tasks")"
        return 0
    fi

    function backup(){
        [[ -e "${confDest}" ]] || wrn "[${scope}] $(inf "No config to export ['${confDest}']")"

        mkdir -p "${confSrcDir}"
         if ! cmp -s "${confDest}" "${confDest}"; then
            sinf "Exporting config ['${confDest}' -> '${confSrc}']"
            cp -v "${confDest}" "${confSrc}" ||\
            smsg "$(wrn "Couldnt backed-up config [${confSrc}]")"
        else
            sinf "No changes — config export is up-to-date 📄✅"
        fi
    }

    function setup(){
        if [[ -e ${confSrc} ]]; then
            sinf "Setting up [${confSrc}]"
            if [[ -d $(dirname "${confSrc}") ]]; then
                mkdir -p "${confDestDir}"
                mkdir -p "${confEnDir}"

                cp -v "${confSrc}" "${confDest}" && \
                ln -svf "${confDest}" "${confEn}"
            else
                wrn "[${scope}] $(inf "Config Dir not found [$(dirname "${confSrc}")]")"
            fi
        else
            wrn "[${scope}] $(inf "Config not found [${confSrc}]")"
        fi
    }
    
    [[ "${_nginxSig:-99}" -eq 0 ]] && return
    case "$1" in
        setup)
            setup
            shift ;;
        -b|b|bkp|backup)
            backup
            shift ;;
        *) echo "Usage: nginx {setup|backup}" ;;
    esac
}

function _db(){
    scope=DB;

    init() {
        :
    }
    patch(){
        sinf "Applying patches"
        # grep -q "SET foreign_key_checks = 0;" "${1:-${dbSrc}.sql}" ||\
        # sed -i '1s/^/SET foreign_key_checks = 0;\n/' "${1:-${dbSrc}.sql}"
        
        sed -i '/\/\*M!999999\\- enable the sandbox mode \*\//d' "${1:-${dbSrc}.sql}"
    }

    clean() {
        local sqlFile="${1:-${dbSrc}.sql}"
        case "$2" in
            std)
                grep -vE '^\s*(--|#|$)|/\*!.*\*/|/\*M!.*\*/' "${sqlFile}"
                return
                shift ;;
            :) echo "Usage: db clean {filename} {std}"
                return
                ;;
        esac

        sinf "Cleaning file '${sqlFile}'"
        sed -E -i \
            -e '/^\s*(--|#)\s*/d' \
            -e '/^\s*\/\*![0-9]+[^*]*\*\/;?\s*$/d' \
            -e '/^\s*\/\*M![^*]*\*\/;?\s*$/d' \
            -e 's@/\*![0-9]+[^*]*\*/;?@@g' \
            -e 's@/\*M![^*]*\*/;?@@g' \
            "${sqlFile}"

    }

    backup(){

        __db_export() {
            mysqldump -h "${dbHost}" -u "${dbName}_u" -p"${dbPass}" "${dbName}" > "${dbSrc}.sql.tmp" ||\
            smsg "$(wrn "Export failed")"
        }

        __db_export_as_root(){
            [ ${UID} -eq 0 ] &&\
            (mysqldump -h "${dbHost}" -u root --skip-password "${dbName}" > "${dbSrc}.sql.tmp") ||\
            smsg "$(wrn "Privileged access required to export '${dbName}'")" && return 1
        }

        mkdir -p "${dbSrcDir}"

        (__db_export || __db_export_as_root) && patch "${dbSrc}.sql.tmp" > /dev/null && {
            if ! cmp -s <(clean "${dbSrc}.sql.tmp" "std") <(clean "${dbSrc}.sql" "std"); then
                sinf "Exporting to ['${dbSrc}.sql']"
                mv -v "${dbSrc}.sql.tmp" "${dbSrc}.sql"
            else
                sinf "No changes — db export is up-to-date 🛢️ ✅"
                rm -f "${dbSrc}.sql.tmp"
            fi
        } || {
            smsg "$(wrn "Export failed")"
            return 1
        }
    }


    setup() {
        sql_files=(
            "${dbSrc}-init.sql"
            "${dbSrc}.sql"
        )

        __db_c() {
            # -create db 
            (mysql -h "${dbHost:-localhost}" -u "${dbName}_u" -p"${dbPass}" -e"
            DROP DATABASE IF EXISTS \`${dbName}\`; 
            CREATE DATABASE IF NOT EXISTS \`${dbName}\`;
            " ) ||\
            (smsg "$(wrn "User '${dbName}_u' wasn't found, creating user for '${dbName}'")"
            return 1)
        }

        __db_c_user() {
            # -create db user
            (mysql -h "${dbHost:-localhost}" -u root --skip-password -e"
                CREATE USER IF NOT EXISTS '${dbName}_u'@'localhost' IDENTIFIED BY '${dbPass}';
                GRANT ALL PRIVILEGES ON \`${dbName}\`.* TO '${dbName}_u'@'localhost';
                FLUSH PRIVILEGES;
            " ) ||\
            (smsg "$(wrn "Failed to create user '${dbName}_u'")"
            return 1)
        }

        # -attempt to create db & user
        (__db_c && sinf "User '${dbName}_u' @ '${dbName}' exists" ) ||\
        [ ${UID} -eq 0 ] && (__db_c_user && __db_c) ||\
        smsg "$(wrn "Priviledge access required to create user '${dbName}_u'")"


        __db_import() {
            mysql -h "${dbHost}" -u "${dbName}_u" -p"${dbPass}" "${dbName}" -e"
                SET FOREIGN_KEY_CHECKS=0; SOURCE "$sql_file"; SET FOREIGN_KEY_CHECKS=1;
            " > /dev/null || return 1
        }
        __db_import_as_root() {
            mysql -h "${dbHost}" -u root --skip-password "${dbName}"  -e"
                SET FOREIGN_KEY_CHECKS=0; SOURCE "$sql_file"; SET FOREIGN_KEY_CHECKS=1;
            " > /dev/null || return 1
        }
        for sql_file in "${sql_files[@]}"; do
            
            if [[ -f "$sql_file" ]]; then
                sinf "Importing : $sql_file"
                
                __db_import ||\
                ([ "$UID" -eq 0 ] && __db_import_as_root) ||\
                smsg "$(wrn "Failed to import : $sql_file")"
            
            else
                smsg "$(wrn "File not found: $sql_file")"
            fi
        done

        sctl_r_mysql=1
    }

    function update(){
        msg "[${scope}]($gitBranch) $(inf "updating - ${dbSrc}.sql")"
        git fetch origin
        git checkout origin/"$gitBranch" -- "${dbSrc}.sql"

        setup
    }

    case "$1" in
        patch)
            patch
            shift ;;
        clean)
            clean
            shift ;;
        clean-std)
            clean "${dbSrc}.sql" "std" 
            shift ;;
        setup|import)
            setup
            shift ;;
        -b|b|bkp|backup|export)
            backup
            shift ;;
        u|update)
            update
            shift ;;
        *) echo "Usage: db {patch|setup|backup}" ;;
    esac
}

function _composer(){
    scope=Composer;

    function backup(){
        :
    }
    function setup(){
        mkdir -p "${composerCacheDir}"
        if [ "$(id -u)" -eq 0 ]; then
            chown -R "${safeUser}":"${safeUser}" "${composerCacheDir}"
            sinf "Switching to user 'www-data' | '${safeUser}'"
            sinf "Installing & Updating"

            sudo -u "${safeUser}" composer install --no-scripts --no-interaction --no-dev --optimize-autoloader ||\
            composer install --no-scripts --no-interaction --no-dev --optimize-autoloader ||\
            smsg "$(wrn "Install failed")"

            sudo -u "${safeUser}" composer update --no-scripts --no-interaction ||\
            composer update --no-scripts --no-interaction ||\
            smsg "$(wrn "Update failed")"
        else
            sinf "Installing & Updating"
            composer install --no-scripts --no-interaction --no-dev --optimize-autoloader ||\
            smsg "$(wrn "Install failed")"
            composer update --no-scripts --no-interaction ||\
            smsg "$(wrn "Update failed")"
        fi
    }

    case "$1" in
        setup) setup
            shift ;;
        -b|b|bkp|backup) backup
            shift ;;
        *) echo "Usage: composer {setup|backup}" ;;
    esac
}

function _laravel(){
    scope=Latavel;

    function backup(){
        :
    }
    function setup(){
        inf "Configuring Laravel..."
        [[ ! -f .env ]] && cp .env.example .env || alrt "Php env file exists, keeping it"

        php artisan key:generate
        php artisan migrate --force
        php artisan storage:link

        inf "Regenerating Laravel config..."
        php artisan config:clear
        php artisan config:cache
        php artisan route:clear
        php artisan view:clear
    }

    case "$1" in
        setup) setup
            shift ;;
        -b|b|bkp|backup) backup
            shift ;;
        *) echo "Usage: larvel {setup|backup}" ;;
    esac
}

function _git(){
    scope=Git;
    _branch="$(msg "(")$(inf "$gitBranch")$(msg ")")\e[0;2m"

    function pull(){
        # smsg "($gitBranch) $(inf "Pulling changes from '$gitBranch'")"
        sinf "${_branch} Pulling changes from '$gitBranch'"

        if git diff --quiet; then
            git pull origin "$gitBranch"
        elif [ "$gitFetchReset" -eq 1 ]; then
            sinf "${_branch} Changes detected, discarding em all"
            git reset --hard HEAD
            git clean -fd
            git pull origin "$gitBranch"
        else
            wrn "Failed to pull changes"
        fi
    }

    function backup(){
        if ! git diff --quiet; then
            local gitMsg="bkp-${gitMsg}"
            sinf "${_branch} Syncing working tree '$gitMsg'"
            git add -A && git commit -m "$gitMsg" && git push
        else
            sinf "${_branch} No changes — working tree is clean, in sync 🕒✅"
        fi
    }

    function sync(){
        local gitMsg="sync-${gitMsg}"
        sinf "${_branch} Staging & commiting changes to '$gitBranch' - $gitMsg"
        git add -A
        git commit -m "$gitMsg" || wrn "Failed to commit changes"

        sinf "${_branch} Fetching changes from '$gitBranch'"
        git pull origin "$gitBranch" --rebase || wrn "Failed to fetch changes"

        sinf "${_branch} Pushing changes to '$gitBranch'"
        git push origin "$gitBranch" || wrn "Failed to push changes"
    }

    gitMsgMeta="$(hostname)-$(date '+%Y-%m-%d-%a')"
    gitMsg="${gitMsg:+$gitMsgMeta $gitMsg}"
    gitMsg="${gitMsg:-$gitMsgMeta:$(date '+%H%M%S')}"

    case "$1" in
        pull)
            shift
            gitFetchReset=1
            pull ;;
        -b|b|bkp|backup)
            shift
            backup ;;
        sync)
            shift
            sync ;;
        *) echo "Usage: git {pull|backup|sync}" ;;
    esac
}

function sctl(){
    local actions=$1
    shift

    for ((i=0; i<${#actions}; i++)); do
        local char="${actions:$i:1}"
        SERVICES="nginx mariadb php${PHP_VER}-fpm.service"
        
        case "$char" in
            e) systemctl enable ${SERVICES} ;;
            r) systemctl restart ${SERVICES} ;;
            x) systemctl stop ${SERVICES} ;;
            d) systemctl disable ${SERVICES} ;;
            s) systemctl status ${SERVICES} ;;
            *) echo "Usage: srvc|service {e(enable)|r(restart)|s(stop)|d(disable)}" ;;
        esac
    done 
}

function _backup(){
    scope=Bkp;

    msg "\n# Backing-up"
    _nginx && _db bkp && _git bkp
}

function _update(){
    scope=Update;

    msg "\n# Perfoming Update"
    fix
    _git pull
    _nginx setup
    _db setup
    #_composer setup
    _laravel setup
}

function _setup(){
    scope=Setup;

    fix
    msg "\n# Setting-up project"
    _pkg setup
    _git pull
    _nginx setup
    _db setup
    _composer setup
    _laravel setup
    sctl er
}

usage() {
    cat <<EOF
Usage: $0 [options]

Options:
  -a, --all          Run all setup tasks (pkg, git pull, nginx, db, composer, laravel, srvc)
  -b, --backup       Perform backup (db and git)
  -u, --update       Perform update (git pull, nginx setup, db setup, laravel setup)
  --pkg              Setup required packages
  --composer         Install/update composer dependencies
  --laravel          Setup Laravel environment
  --srvc, --service  Enable, disable, restart, or stop services (nginx, mariadb, php)
  --git              Run git commands (pull, sync, backup)
  --db               Setup, backup, or patch the database
  -m, --message      Commit message for git backup/sync
  -h, --help         Display this help message
EOF
}

# Call usage function in case of invalid argument or help request
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    usage
    exit 0
fi

gitBranch="${gitBranch:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'main')}"

while [[ $# -gt 0 ]]; do
    case $1 in
        -a|a|--all|all|--setup|setup)
            shift
            _setup
            ;;
        -b|b|--bkp|bkp|--backup|backup)
            shift
            bkpSig=1
            ;;
        -u|u|--update|update)
            shift
            updateSig=1
            ;;
        -f|f|--fix|fix)
            fix
            shift
            ;;
        --pkg|pkg)
            _pkg $2
            shift 2
            ;;
        --nginx|nginx)
            _nginx "$2"
            shift 2
            ;;
        --no-nginx|no-nginx|--nonginx|nonginx)
            nginxSig=0;
            shift
            ;;
        --composer|composer)
            _composer "$2"
            shift 2
            ;;
        --laravel|laravel)
            _laravel "$2"
            shift 2
            ;;
        --srvc|--sctl|sctl|srvc|service|services)
            sctl "$2"
            shift 2
            ;;
        --git|git)
            git_arg="$2"
            gitSig=1
            shift 2
            ;;
        --branch|branch)
            gitBranch="$2"
            shift 2
            ;;
        -m|m|--message|message)
            gitMsg="$2"
            shift 2
            ;;
        --db|db)
            db_arg="$2"
            # _db "$@"
            dbSig=1
            shift 2
            ;;
        --db-clean-std|db-clean-std)
            _db clean-std
            shift
            ;;
        --dbname|dbname)
            dbName="$2"
            d_db_arg=setup
            dbSig=1
            shift 2
            ;;
        -h|help)
            usage
            exit 0
            ;;
        *)
            wrn "Error: Invalid option $1"
            usage
            exit 1
            ;;
    esac
done


[ ${gitSig:-0} -eq 1 ] && _git "$git_arg" || :
[ ${dbSig:-0} -eq 1 ] && _db "${db_arg:-$d_db_arg}" || :
[ ${bkpSig:-0} -eq 1 ] && _backup || :
[ ${updateSig:-0} -eq 1 ] && _update || :
[ ${fixSig:-0} -eq 1 ] && fix || :
# [ $? -eq 0 ] && inf "Setup complete.." || wrn "Setup failed!!"
