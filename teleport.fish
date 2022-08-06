function teleport 
    set -x OUTPUT (gotp-exe $argv)

    if test $status = 2
        cd "$OUTPUT"
    else
        printf "$(echo $OUTPUT)"
    end
end