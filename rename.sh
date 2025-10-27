#!/bin/bash
find ./templates -type f ! -name "*.tmpl" -exec mv "{}" "{}.tmpl" \;
