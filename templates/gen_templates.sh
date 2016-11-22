#!/usr/bin/env bash
go-bindata-assetfs -pkg templates -nocompress -prefix templates -ignore .+\.go$ -ignore .+\.sh$ .
mv bindata_assetfs.go templates.go

