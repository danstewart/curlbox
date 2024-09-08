#!/usr/bin/env bash

function main() {
    action="$1"
    shift

    if [[ $action == "create" ]]; then
        _debug "action=create"
        create $@
    elif [[ $action == "run" ]]; then
        _debug "action=run"
        run $@
    else
        echo "Usage: "
        echo "    $0 create path/to/new/curlbox"
        echo "    $0 run path/to/curlbox/script-to-run"
        exit 1
    fi
}

function create() {
    box_location="$1"

    if [[ -z $box_location ]]; then
    echo "Usage: $0 <box_location>"
    echo ""
    echo "Creates a new curlbox in the given directory"
    echo "If the directory does not exist, it will be created"
    exit 1
    fi

    mkdir -p "$box_location"
    touch "$box_location/.curlbox-root"

    echo "Created curlbox $box_location"
}

function run() {
    what_to_run="$1"

    if [[ ! -f "$what_to_run" ]]; then
        echo "File '$what_to_run' not found"
        exit 1
    fi

    # Load all variable files
    # Loads from outer directory down into the script directory
    # Loads _vars and then _secret_vars in each directory
    function load_vars() {
        script_dir="$( cd "$( dirname "${what_to_run}" )" && pwd )"

        _debug "what_to_run=$what_to_run"
        _debug "script_dir=$script_dir"

        var_files=()
        dir=$script_dir

        # Keep going up the directory tree until we find the root of the curlbox
        # This is the directory that contains the .curlbox-root file
        while true; do
            if [[ -f "${dir}/_secret_vars" ]]; then
            var_files+=("${dir}/_secret_vars")
            fi

            if [[ -f "${dir}/_vars" ]]; then
            var_files+=("${dir}/_vars")
            fi

            if [[ $dir == "/" || -f "$dir/.curlbox-root" ]]; then
            break
            fi

            dir=$(realpath $dir/..)
        done

        declare -a sorted_var_files
        _reverse_array var_files sorted_var_files
        _debug "sorted_var_files=${sorted_var_files[*]}"

        set -o allexport

        if [[ -f "${root_dir}/_global_vars" ]]; then
            _debug "Loading global vars from ${root_dir}/_global_vars"
            source ${root_dir}/_global_vars
        fi

        for var_file in ${sorted_var_files[@]}; do
            _debug "Loading vars from $var_file"
            source $var_file
        done

        set +o allexport
    }

    load_vars
    shift  # Pop off the script name and pass all other args through to the script

    # Run the script
    $what_to_run $@
}

function _debug() {
    if [[ -n $DEBUG && $DEBUG=1 ]]; then
        echo "[DEBUG] $1"
    fi
}

_reverse_array() {
    # first argument is the array to reverse
    # second is the output array
    declare -n arr="$1" rev="$2"
    for i in "${arr[@]}"
    do
        rev=("$i" "${rev[@]}")
    done
}

main $@
