#!/bin/bash
find ./templates -type f -name "*.tmpl" -exec sed -i '' 's/backend-example/{{ .Module }}/g' {} +
