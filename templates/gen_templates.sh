#!/usr/bin/env bash
go-bindata-assetfs -pkg templates -nocompress -prefix templates -ignore .+\.go$ -debug .
mv bindata_assetfs.go templates.go

