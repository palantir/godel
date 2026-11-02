#!/usr/bin/env bash
export CIRCLECI_TEMPLATE=go-distribution-oss
export ADDITIONAL_MODULES=pkg/products/v2
export ENABLE_INTEGRATION_TEST=true
export USE_LEGACY_GO_TOOLCHAIN_SELECTION=false
