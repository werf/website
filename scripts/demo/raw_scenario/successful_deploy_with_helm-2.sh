helm upgrade --wait --install --create-namespace --namespace sample-app sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app.test
curl sample-app.test/ping
