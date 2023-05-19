helm upgrade --wait --install --create-namespace --namespace sample-app-2 sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app-2.test
curl sample-app-2.test/ping
