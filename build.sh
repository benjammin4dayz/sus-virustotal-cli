#!/bin/bash

entry_file="main.go"
output_dir="dist"
os_list=("linux" "windows" "darwin")
arch_list=("amd64" "arm" "arm64")

rm -rf "$output_dir" && mkdir -p "$output_dir"

for os in "${os_list[@]}"; do
    for arch in "${arch_list[@]}"; do
        if [ "$os" = "darwin" ] && [ "$arch" = "arm" ]; then continue; fi
        
        output_file="sus-virustotal-cli-${os}_${arch}/sus"
        case $os in
            windows)
                output_file="$output_file.exe"
            ;;
            darwin|linux)
                chmod +x "$output_dir/$output_file" 2>/dev/null
            ;;
        esac
        GOOS="$os" GOARCH="$arch" go build -o "$output_dir/$output_file" $entry_file
        echo "$output_file: $(sha256sum "$output_dir/$output_file" | awk '{print $1}')"
    done
done

for dir in "$output_dir"/*; do
    case $dir in
        *linux*|*darwin*)
            tar -czf "${dir}.tar.gz" -C "$dir" .
        ;;
        *windows*)
            zip "${dir}.zip" -j "$dir"/*
        ;;
    esac
done
