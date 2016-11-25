#!/usr/bin/env bash
go-bindata-assetfs -pkg templates -prefix templates -ignore .+\.go$ -ignore .+\.sh$ .
mv bindata_assetfs.go templates.go

